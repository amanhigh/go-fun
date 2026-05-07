import { createAsyncFeedbackState, runAsyncFeedback } from '../../../lib/async_feedback';
import { formatDateInputValue } from '../../../lib/date';
import { normalizeTag } from '../../../lib/tags';
import type { Journal, JournalUpdate, JournalUpdateRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function localToday(): string {
	return formatDateInputValue(new Date());
}

// ===== Quick Action Rules =====
// Only TAKEN and REJECTED journals qualify for a quick action.

const QUICK_ACTION_MAP: Record<string, { status: string; label: string; className: string }> = {
	TAKEN: { status: 'JUST_LOSS', label: 'Mark Just Loss', className: 'journal-quick-status-loss' },
	REJECTED: { status: 'BROKEN', label: 'Mark Broken', className: 'journal-quick-status-broken' },
};

function quickActionFor(journal: Journal) {
	const journalType = normalizeTag(journal.type ?? '');
	const action = QUICK_ACTION_MAP[journalType];
	if (!action) return null;
	return normalizeTag(journal.status) === action.status ? null : action;
}

// ===== Update Planning =====

function toggleReviewUpdate(journal: Journal): { payload: JournalUpdateRequest; message: string } {
	const reviewedAt = journal.reviewed_at ? null : localToday();
	return {
		payload: { reviewed_at: reviewedAt },
		message: reviewedAt ? 'Journal marked reviewed.' : 'Journal marked not reviewed.',
	};
}

function quickStatusReviewUpdate(
	status: string,
	label: string,
): { payload: JournalUpdateRequest; message: string } {
	return {
		payload: { status, reviewed_at: localToday() },
		message: `${label} applied and journal marked reviewed.`,
	};
}

function applyReviewUpdate(journal: Journal, update: JournalUpdate): void {
	journal.reviewed_at = update.reviewed_at;
	journal.status = update.status;
}

// ===== Action Context =====

function reviewActionContext(pg: JournalDetailPageProvider) {
	const { current, client } = pg();
	const { journal, journalId } = current;
	if (!journal) return null;
	return { current, client, journal, journalId };
}

// ===== Exported Concern =====

export function NewReviewActionsConcern(pg: JournalDetailPageProvider) {
	return {
		...createAsyncFeedbackState('submitting', 'message', 'messageType'),

		get feedbackClass(): string {
			return this.messageType === 'success' ? 'journal-feedback-success' : 'journal-feedback-error';
		},

		toggleLabel() {
			return pg().current.journal?.reviewed_at ? 'Mark Pending' : 'Mark Reviewed';
		},
		buttonClass() {
			return pg().current.journal?.reviewed_at
				? 'journal-review-toggle-pending'
				: 'journal-review-toggle-reviewed';
		},
		quickAction() {
			const journal = pg().current.journal;
			if (!journal) return null;
			return quickActionFor(journal);
		},

		async toggle() {
			const ctx = reviewActionContext(pg);
			if (!ctx) return;
			await runAsyncFeedback(this, async () => {
				const { payload, message } = toggleReviewUpdate(ctx.journal);
				const envelope = await ctx.client.updateReview(ctx.journalId, payload);
				applyReviewUpdate(ctx.journal, envelope.data);
				this.message = message;
				this.messageType = 'success';
			}, 'Unable to update review date.');
		},

		async applyQuickStatus() {
			const ctx = reviewActionContext(pg);
			if (!ctx) return;
			const action = quickActionFor(ctx.journal);
			if (!action) return;
			await runAsyncFeedback(this, async () => {
				const { payload, message } = quickStatusReviewUpdate(action.status, action.label);
				const envelope = await ctx.client.updateReview(ctx.journalId, payload);
				applyReviewUpdate(ctx.journal, envelope.data);
				this.message = message;
				this.messageType = 'success';
				await pg().sidebar.reviewQueue.load();
			}, 'Unable to update journal status.');
		},
	};
}

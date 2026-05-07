import { createAsyncFeedbackState, runAsyncFeedback } from '../../../lib/async_feedback';
import type { FeedbackState } from '../../../lib/async_feedback';
import { formatDateInputValue } from '../../../lib/date';
import { normalizeTag } from '../../../lib/tags';
import { createQuickAction } from '../../../lib/quick_action';
import type { DisplaySpec } from '../../../types/presentation_concern';
import type { QuickAction } from '../../../types/quick_action';
import type { Journal, JournalUpdate, JournalUpdateRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function localToday(): string {
	return formatDateInputValue(new Date());
}

// ===== Quick Action Maps =====

const STATUS_ACTIONS: Record<string, { status: string; label: string; className: string }> = {
	TAKEN: { status: 'JUST_LOSS', label: 'Mark Just Loss', className: 'journal-quick-status-loss' },
	REJECTED: { status: 'BROKEN', label: 'Mark Broken', className: 'journal-quick-status-broken' },
};

function applyReviewUpdate(journal: Journal, update: JournalUpdate): void {
	journal.reviewed_at = update.reviewed_at;
	journal.status = update.status;
}

// ===== Display Helpers (always return a ready spec) =====

function reviewToggleDisplay(journal: Journal): DisplaySpec {
	const hasReview = !!journal.reviewed_at;
	return {
		icon: '',
		text: hasReview ? 'Mark Pending' : 'Mark Reviewed',
		class: hasReview ? 'journal-review-toggle-pending' : 'journal-review-toggle-reviewed',
	};
}

function reviewStatusDisplay(journal: Journal): DisplaySpec {
	const journalType = normalizeTag(journal.type ?? '');
	const action = STATUS_ACTIONS[journalType];
	return action
		? { icon: '', text: action.label, class: action.className }
		: { icon: '', text: '', class: '' };
}

// ===== Active Checks =====

function reviewStatusActive(journal: Journal): boolean {
	const journalType = normalizeTag(journal.type ?? '');
	const action = STATUS_ACTIONS[journalType];
	return !!action && normalizeTag(journal.status) !== action.status;
}

// ===== Async Action Handlers =====

async function toggleReview(feedback: FeedbackState, pg: JournalDetailPageProvider): Promise<void> {
	const { current, client } = pg();
	const { journal, journalId } = current;
	if (!journal) return;
	await runAsyncFeedback(feedback, async () => {
		const reviewedAt = journal.reviewed_at ? null : localToday();
		const payload: JournalUpdateRequest = { reviewed_at: reviewedAt };
		const envelope = await client.updateReview(journalId, payload);
		applyReviewUpdate(journal, envelope.data);
		feedback.message = reviewedAt ? 'Journal marked reviewed.' : 'Journal marked not reviewed.';
		feedback.messageType = 'success';
	}, 'Unable to update review date.');
}

async function applyReviewStatus(feedback: FeedbackState, pg: JournalDetailPageProvider): Promise<void> {
	const { current, client } = pg();
	const { journal, journalId } = current;
	if (!journal) return;
	const journalType = normalizeTag(journal.type ?? '');
	const action = STATUS_ACTIONS[journalType];
	if (!action) return;
	if (normalizeTag(journal.status) === action.status) return;
	await runAsyncFeedback(feedback, async () => {
		const payload: JournalUpdateRequest = { status: action.status, reviewed_at: localToday() };
		const envelope = await client.updateReview(journalId, payload);
		applyReviewUpdate(journal, envelope.data);
		feedback.message = `${action.label} applied and journal marked reviewed.`;
		feedback.messageType = 'success';
		await pg().sidebar.reviewQueue.load();
	}, 'Unable to update journal status.');
}

// ===== Exported Concern =====

export function NewReviewActionsConcern(pg: JournalDetailPageProvider) {
	return {
		...createAsyncFeedbackState('submitting', 'message', 'messageType'),

		get feedbackClass(): string {
			return this.messageType === 'success' ? 'journal-feedback-success' : 'journal-feedback-error';
		},

		isSubmitting() {
			return this.submitting;
		},

		actions(): QuickAction[] {
			const journal = pg().current.journal;
			if (!journal) return [];

			return [
				createQuickAction('review-toggle', () => true, reviewToggleDisplay(journal), () => toggleReview(this, pg)),
				createQuickAction('review-status', () => reviewStatusActive(journal), reviewStatusDisplay(journal), () => applyReviewStatus(this, pg)),
			].filter((action) => action.isActive());
		},
	};
}

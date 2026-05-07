import { createAsyncFeedbackState, runAsyncFeedback } from '../../../lib/async_feedback';
import { formatDateInputValue } from '../../../lib/date';
import { normalizeTag } from '../../../lib/tags';
import type { DisplaySpec } from '../../../types/presentation_concern';
import type { QuickAction } from '../../../types/quick_action';
import type { Journal, JournalUpdate, JournalUpdateRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function localToday(): string {
	return formatDateInputValue(new Date());
}

// ===== Quick Action Maps =====

const QUICK_ACTION_MAP: Record<string, { status: string; label: string; className: string }> = {
	TAKEN: { status: 'JUST_LOSS', label: 'Mark Just Loss', className: 'journal-quick-status-loss' },
	REJECTED: { status: 'BROKEN', label: 'Mark Broken', className: 'journal-quick-status-broken' },
};

function applyReviewUpdate(journal: Journal, update: JournalUpdate): void {
	journal.reviewed_at = update.reviewed_at;
	journal.status = update.status;
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
			const toggle: QuickAction = {
				id: 'review-toggle',
				display: () => {
					if (!pg().current.journal) return { icon: '', text: '', class: '' };
					const hasReview = !!pg().current.journal.reviewed_at;
					return {
						icon: '',
						text: hasReview ? 'Mark Pending' : 'Mark Reviewed',
						class: hasReview ? 'journal-review-toggle-pending' : 'journal-review-toggle-reviewed',
					};
				},
				apply: async () => {
					const { current, client } = pg();
					const { journal, journalId } = current;
					if (!journal) return;
					await runAsyncFeedback(this, async () => {
						const reviewedAt = journal.reviewed_at ? null : localToday();
						const payload: JournalUpdateRequest = { reviewed_at: reviewedAt };
						const envelope = await client.updateReview(journalId, payload);
						applyReviewUpdate(journal, envelope.data);
						this.message = reviewedAt ? 'Journal marked reviewed.' : 'Journal marked not reviewed.';
						this.messageType = 'success';
					}, 'Unable to update review date.');
				},
			};

			const status: QuickAction = {
				id: 'review-status',
				display: () => {
					const journal = pg().current.journal;
					if (!journal) return { icon: '', text: '', class: '' };
					const journalType = normalizeTag(journal.type ?? '');
					const action = QUICK_ACTION_MAP[journalType];
					if (!action) return { icon: '', text: '', class: '' };
					if (normalizeTag(journal.status) === action.status) return { icon: '', text: '', class: '' };
					return { icon: '', text: action.label, class: action.className };
				},
				apply: async () => {
					const { current, client } = pg();
					const { journal, journalId } = current;
					if (!journal) return;
					const journalType = normalizeTag(journal.type ?? '');
					const action = QUICK_ACTION_MAP[journalType];
					if (!action) return;
					if (normalizeTag(journal.status) === action.status) return;
					await runAsyncFeedback(this, async () => {
						const payload: JournalUpdateRequest = { status: action.status, reviewed_at: localToday() };
						const envelope = await client.updateReview(journalId, payload);
						applyReviewUpdate(journal, envelope.data);
						this.message = `${action.label} applied and journal marked reviewed.`;
						this.messageType = 'success';
						await pg().sidebar.reviewQueue.load();
					}, 'Unable to update journal status.');
				},
			};

			return [toggle, status];
		},
	};
}

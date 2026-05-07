import { createAsyncFeedbackState } from '../../../shared/async_feedback';
import { formatDateInputValue } from '../../../shared/date';
import { getErrorMessage } from '../../../shared/error';
import { normalizeTag } from '../../../shared/tags';
import type { JournalUpdateRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function localToday(): string {
	return formatDateInputValue(new Date());
}

/** Shared sentinel returned when no quick action is applicable. */
const NO_QUICK_ACTION = { status: '', label: '', className: '' };

/**
 * Map from journal type to the quick action definition.
 * Only TAKEN and REJECTED journals qualify for a quick action.
 */
const QUICK_ACTION_MAP: Record<string, { status: string; label: string; className: string }> = {
	TAKEN: { status: 'JUST_LOSS', label: 'Mark Just Loss', className: 'journal-quick-status-loss' },
	REJECTED: { status: 'BROKEN', label: 'Mark Broken', className: 'journal-quick-status-broken' },
};

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
			if (!journal) return NO_QUICK_ACTION;
			const journalType = normalizeTag(journal.type ?? '');
			const action = QUICK_ACTION_MAP[journalType];
			if (!action) return NO_QUICK_ACTION;
			return normalizeTag(journal.status) === action.status ? NO_QUICK_ACTION : action;
		},

		async toggle() {
			const { current, client } = pg();
			const { journal, journalId } = current;
			if (!journal || this.submitting) return;
			this.submitting = true;
			this.message = '';
			this.messageType = 'error';
			try {
				const markingReviewed = !journal.reviewed_at;
				const reviewedAt = markingReviewed ? localToday() : null;
				const payload: JournalUpdateRequest = { reviewed_at: reviewedAt };
				const envelope = await client.updateReview(journalId, payload);
				const updatedJournal = current.journal;
				if (updatedJournal) {
					updatedJournal.reviewed_at = envelope.data.reviewed_at;
				}
				this.messageType = 'success';
				this.message = markingReviewed ? 'Journal marked reviewed.' : 'Journal marked not reviewed.';
			} catch (err) {
				this.message = getErrorMessage(err, 'Unable to update review date.');
				this.messageType = 'error';
			} finally {
				this.submitting = false;
			}
		},

		async applyQuickStatus() {
			const { current, client } = pg();
			const { journalId } = current;
			const action = this.quickAction();
			if (this.submitting || !action.status) return;
			const { status, label } = action;
			const payload: JournalUpdateRequest = { status, reviewed_at: localToday() };
			this.submitting = true;
			this.message = '';
			this.messageType = 'error';
			try {
				const envelope = await client.updateReview(journalId, payload);
				const updatedJournal = current.journal;
				if (updatedJournal) {
					updatedJournal.status = envelope.data.status;
					updatedJournal.reviewed_at = envelope.data.reviewed_at;
				}
				this.messageType = 'success';
				this.message = `${label} applied and journal marked reviewed.`;
				await pg().sidebar.reviewQueue.load();
			} catch (err) {
				this.message = getErrorMessage(err, 'Unable to update journal status.');
				this.messageType = 'error';
			} finally {
				this.submitting = false;
			}
		},
	};
}

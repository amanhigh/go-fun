import { createAsyncFeedbackState } from '../../../shared/async_feedback';
import { formatDateInputValue } from '../../../shared/date';
import { getErrorMessage } from '../../../shared/error';
import { normalizeTag } from '../../../shared/tags';
import type { JournalUpdateRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function localToday(): string {
	return formatDateInputValue(new Date());
}

export function createReviewActionsState() {
	return createAsyncFeedbackState('submitting', 'message', 'messageType');
}

export function NewReviewActionsConcern(pg: JournalDetailPageProvider) {
	return {
		...createReviewActionsState(),

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
			const journalType = normalizeTag(journal?.type ?? '');
			if (journalType !== 'TAKEN' && journalType !== 'REJECTED') return null;
			const status = journalType === 'TAKEN' ? 'JUST_LOSS' : 'BROKEN';
			const label = status === 'JUST_LOSS' ? 'Mark Just Loss' : 'Mark Broken';
			const className = status === 'JUST_LOSS'
				? 'journal-quick-status-loss'
				: 'journal-quick-status-broken';
			if (!journal) return null;
			const isActive = normalizeTag(journal.status) === status;
			return isActive ? null : { status, label, className };
		},
		hasQuickAction() {
			return this.quickAction() !== null;
		},
		quickLabel() {
			return this.quickAction()?.label ?? 'Update Status';
		},
		quickButtonClass() {
			return this.quickAction()?.className ?? '';
		},

		async toggle() {
			const journal = pg().current.journal;
			if (!journal || this.submitting) return;
			this.submitting = true;
			this.message = '';
			this.messageType = 'error';
			try {
				const reviewedAt = journal.reviewed_at ? null : localToday();
				const payload: JournalUpdateRequest = { reviewed_at: reviewedAt };
				const envelope = await pg().client.updateReview(pg().current.journalId, payload);
				const current = pg().current.journal;
				if (current) {
					current.reviewed_at = envelope.data.reviewed_at;
				}
				this.messageType = 'success';
				this.message = reviewedAt ? 'Journal marked reviewed.' : 'Journal marked not reviewed.';
			} catch (err) {
				this.message = getErrorMessage(err, 'Unable to update review date.');
				this.messageType = 'error';
			} finally {
				this.submitting = false;
			}
		},

		async applyQuickStatus() {
			const journal = pg().current.journal;
			const action = this.quickAction();
			if (!journal || this.submitting || !action) return;
			const { status } = action;
			this.submitting = true;
			this.message = '';
			this.messageType = 'error';
			try {
				const envelope = await pg().client.updateReview(pg().current.journalId, { status, reviewed_at: localToday() });
				const current = pg().current.journal;
				if (current) {
					current.status = envelope.data.status;
					current.reviewed_at = envelope.data.reviewed_at;
				}
				this.messageType = 'success';
				this.message = `${action.label} applied and journal marked reviewed.`;
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

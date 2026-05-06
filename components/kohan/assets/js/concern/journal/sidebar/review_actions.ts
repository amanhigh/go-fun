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
			return this.messageType === 'success' ? 'text-emerald-700' : 'text-rose-700';
		},

		toggleLabel(this: any) {
			return pg().current.journal?.reviewed_at ? 'Mark Pending' : 'Mark Reviewed';
		},
		buttonClass(this: any) {
			return pg().current.journal?.reviewed_at
				? 'border-amber-300 bg-amber-50 text-amber-800 hover:bg-amber-100 focus:border-amber-400 focus:ring-amber-200'
				: 'border-emerald-300 bg-emerald-50 text-emerald-800 hover:bg-emerald-100 focus:border-emerald-400 focus:ring-emerald-200';
		},
		quickStatus(this: any) {
			const journalType = normalizeTag(pg().current.journal?.type ?? '');
			if (journalType === 'TAKEN') return 'JUST_LOSS';
			if (journalType === 'REJECTED') return 'BROKEN';
			return '';
		},
		quickLabel(this: any) {
			const status = this.quickStatus();
			if (status === 'JUST_LOSS') return 'Mark Just Loss';
			if (status === 'BROKEN') return 'Mark Broken';
			return 'Update Status';
		},
		hasQuickAction(this: any) {
			const targetStatus = this.quickStatus();
			const journal = pg().current.journal;
			if (!targetStatus || !journal) return false;
			return normalizeTag(journal.status) !== targetStatus;
		},
		quickButtonClass(this: any) {
			return this.quickStatus() === 'JUST_LOSS'
				? 'border-rose-300 bg-rose-50 text-rose-800 hover:bg-rose-100 focus:border-rose-400 focus:ring-rose-200'
				: 'border-violet-300 bg-violet-50 text-violet-800 hover:bg-violet-100 focus:border-violet-400 focus:ring-violet-200';
		},

		async toggle(this: any) {
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

		async applyQuickStatus(this: any) {
			const journal = pg().current.journal;
			if (!journal || this.submitting || !this.hasQuickAction()) return;
			const status = this.quickStatus();
			if (!status) return;
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
				this.message = `${this.quickLabel()} applied and journal marked reviewed.`;
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

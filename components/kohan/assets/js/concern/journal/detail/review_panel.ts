import { createAsyncFeedbackState } from '../../../shared/async_feedback';
import { formatDateInputValue } from '../../../shared/date';
import { getErrorMessage } from '../../../shared/error';
import { normalizeTag } from '../../../shared/tags';
import type { JournalUpdateRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider, ReviewState } from '../../../types/journal_detail_concern';

function localToday(): string {
	return formatDateInputValue(new Date());
}

export function createReviewState(): ReviewState {
	return {
		...createAsyncFeedbackState('reviewSubmitting', 'reviewMessage', 'reviewMessageType'),
		reviewQueue: [],
		reviewQueueLoading: false,
		reviewQueueError: '',
	};
}

export function NewReviewConcern(pg: JournalDetailPageProvider) {
	return {
		reviewToggleLabel(this: any) {
			return pg().journal?.reviewed_at ? 'Mark Pending' : 'Mark Reviewed';
		},
		reviewButtonClass(this: any) {
			return pg().journal?.reviewed_at
				? 'border-amber-300 bg-amber-50 text-amber-800 hover:bg-amber-100 focus:border-amber-400 focus:ring-amber-200'
				: 'border-emerald-300 bg-emerald-50 text-emerald-800 hover:bg-emerald-100 focus:border-emerald-400 focus:ring-emerald-200';
		},
		quickReviewStatus(this: any) {
			const journalType = normalizeTag(pg().journal?.type ?? '');
			if (journalType === 'TAKEN') return 'JUST_LOSS';
			if (journalType === 'REJECTED') return 'BROKEN';
			return '';
		},
		quickReviewLabel(this: any) {
			const status = this.quickReviewStatus();
			if (status === 'JUST_LOSS') return 'Mark Just Loss';
			if (status === 'BROKEN') return 'Mark Broken';
			return 'Update Status';
		},
		hasQuickReviewAction(this: any) {
			const targetStatus = this.quickReviewStatus();
			const journal = pg().journal;
			if (!targetStatus || !journal) return false;
			return normalizeTag(journal.status) !== targetStatus;
		},
		quickReviewButtonClass(this: any) {
			return this.quickReviewStatus() === 'JUST_LOSS'
				? 'border-rose-300 bg-rose-50 text-rose-800 hover:bg-rose-100 focus:border-rose-400 focus:ring-rose-200'
				: 'border-violet-300 bg-violet-50 text-violet-800 hover:bg-violet-100 focus:border-violet-400 focus:ring-violet-200';
		},
		applyReviewUpdate(this: any, payload: JournalUpdateRequest, successMessage: string) {
			return (async () => {
				const envelope = await pg().client.updateReview(pg().journalId, payload);
				const journal = pg().journal;
				if (journal) {
					journal.status = envelope.data.status;
					journal.reviewed_at = envelope.data.reviewed_at;
				}
				this.reviewMessageType = 'success';
				this.reviewMessage = successMessage;
				await this.loadReviewQueue();
			})();
		},
		async toggleReview(this: any) {
			const journal = pg().journal;
			if (!journal || this.reviewSubmitting) return;
			this.reviewSubmitting = true;
			this.reviewMessage = '';
			this.reviewMessageType = 'error';
			try {
				const reviewedAt = journal.reviewed_at ? null : localToday();
				const payload: JournalUpdateRequest = { reviewed_at: reviewedAt };
				await this.applyReviewUpdate(
					payload,
					reviewedAt ? 'Journal marked reviewed.' : 'Journal marked not reviewed.',
				);
			} catch (err) {
				this.reviewMessage = getErrorMessage(err, 'Unable to update review date.');
				this.reviewMessageType = 'error';
			} finally {
				this.reviewSubmitting = false;
			}
		},
		async applyQuickReviewStatus(this: any) {
			if (!pg().journal || this.reviewSubmitting || !this.hasQuickReviewAction()) return;
			// journal guard is in hasQuickReviewAction, safe to use optional chaining below
			const status = this.quickReviewStatus();
			if (!status) return;
			this.reviewSubmitting = true;
			this.reviewMessage = '';
			this.reviewMessageType = 'error';
			try {
				await this.applyReviewUpdate(
					{ status, reviewed_at: localToday() },
					`${this.quickReviewLabel()} applied and journal marked reviewed.`,
				);
			} catch (err) {
				this.reviewMessage = getErrorMessage(err, 'Unable to update journal status.');
				this.reviewMessageType = 'error';
			} finally {
				this.reviewSubmitting = false;
			}
		},
		async loadReviewQueue(this: any) {
			this.reviewQueueLoading = true;
			this.reviewQueueError = '';
			try {
				const envelope = await pg().client.list(0, 10, { reviewed: 'false', sortBy: 'created_at', sortOrder: 'asc' });
				this.reviewQueue = envelope.data?.journals ?? [];
			} catch (err) {
				this.reviewQueueError = getErrorMessage(err, 'Unable to load review queue.');
			} finally {
				this.reviewQueueLoading = false;
			}
		},
	};
}

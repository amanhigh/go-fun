import type {
	Envelope,
	Journal,
	JournalList,
	JournalReviewStatusResponse,
	JournalReviewUpdate,
} from './journal_models';
import { normalizeTag } from './journal_detail_formatters';

function localToday(): string {
	const today = new Date();
	const year = today.getFullYear();
	const month = `${today.getMonth() + 1}`.padStart(2, '0');
	const day = `${today.getDate()}`.padStart(2, '0');
	return `${year}-${month}-${day}`;
}

export function createJournalDetailReview() {
	return {
		reviewToggleLabel(this: any) {
			return this.journal?.reviewed_at ? 'Mark Pending' : 'Mark Reviewed';
		},
		reviewButtonClass(this: any) {
			return this.journal?.reviewed_at
				? 'border-amber-300 bg-amber-50 text-amber-800 hover:bg-amber-100 focus:border-amber-400 focus:ring-amber-200'
				: 'border-emerald-300 bg-emerald-50 text-emerald-800 hover:bg-emerald-100 focus:border-emerald-400 focus:ring-emerald-200';
		},
		quickReviewStatus(this: any) {
			const journalType = normalizeTag(this.journal?.type ?? '');
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
			if (!targetStatus || !this.journal) return false;
			return normalizeTag(this.journal.status) !== targetStatus;
		},
		quickReviewButtonClass(this: any) {
			return this.quickReviewStatus() === 'JUST_LOSS'
				? 'border-rose-300 bg-rose-50 text-rose-800 hover:bg-rose-100 focus:border-rose-400 focus:ring-rose-200'
				: 'border-violet-300 bg-violet-50 text-violet-800 hover:bg-violet-100 focus:border-violet-400 focus:ring-violet-200';
		},
		applyReviewUpdate(this: any, payload: JournalReviewUpdate, successMessage: string, errorMessage: string) {
			return (async () => {
				const response = await fetch(`/v1/api/journals/${this.journalId}`, {
					method: 'PATCH',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(payload),
				});
				if (!response.ok) throw new Error(response.status === 404 ? 'Journal not found' : errorMessage);
				const envelope = (await response.json()) as Envelope<JournalReviewStatusResponse>;
				if (this.journal) {
					this.journal.status = envelope.data.status;
					this.journal.reviewed_at = envelope.data.reviewed_at;
				}
				this.reviewMessageType = 'success';
				this.reviewMessage = successMessage;
				await this.loadReviewQueue();
			})();
		},
		async toggleReview(this: any) {
			if (!this.journal || this.reviewSubmitting) return;
			this.reviewSubmitting = true;
			this.reviewMessage = '';
			this.reviewMessageType = 'error';
			try {
				const reviewedAt = this.journal.reviewed_at ? null : localToday();
				const payload: JournalReviewUpdate = { reviewed_at: reviewedAt };
				await this.applyReviewUpdate(
					payload,
					reviewedAt ? 'Journal marked reviewed.' : 'Journal marked not reviewed.',
					'Failed to update review date',
				);
			} catch (err) {
				this.reviewMessage = err instanceof Error ? err.message : 'Unable to update review date.';
				this.reviewMessageType = 'error';
			} finally {
				this.reviewSubmitting = false;
			}
		},
		async applyQuickReviewStatus(this: any) {
			if (!this.journal || this.reviewSubmitting || !this.hasQuickReviewAction()) return;
			const status = this.quickReviewStatus();
			if (!status) return;
			this.reviewSubmitting = true;
			this.reviewMessage = '';
			this.reviewMessageType = 'error';
			try {
				await this.applyReviewUpdate(
					{ status, reviewed_at: localToday() },
					`${this.quickReviewLabel()} applied and journal marked reviewed.`,
					'Failed to update journal status',
				);
			} catch (err) {
				this.reviewMessage = err instanceof Error ? err.message : 'Unable to update journal status.';
				this.reviewMessageType = 'error';
			} finally {
				this.reviewSubmitting = false;
			}
		},
		async loadReviewQueue(this: any) {
			this.reviewQueueLoading = true;
			this.reviewQueueError = '';
			try {
				const response = await fetch('/v1/api/journals?reviewed=false&sort-by=created_at&sort-order=asc&limit=10');
				if (!response.ok) throw new Error('Failed to load review queue');
				const envelope = (await response.json()) as Envelope<JournalList>;
				this.reviewQueue = envelope.data?.journals ?? [] as Journal[];
			} catch (err) {
				this.reviewQueueError = err instanceof Error ? err.message : 'Unable to load review queue.';
			} finally {
				this.reviewQueueLoading = false;
			}
		},
	};
}

import type { JournalClient } from '../client/journal';
import { getErrorMessage } from '../shared/error';

type CreateJournalDetailPageActionsInput = {
	journalClient: JournalClient;
};

export function createJournalDetailPageActions({ journalClient }: CreateJournalDetailPageActionsInput) {
	return {
		async loadJournal(this: any) {
			this.loading = true;
			this.errorMessage = '';
			try {
				const envelope = await journalClient.get(this.journalId);
				const journal = envelope.data;
				this.journal = journal
					? {
						...journal,
						images: journal.images ?? [],
						tags: journal.tags ?? [],
						notes: [...(journal.notes ?? [])].sort((left, right) => {
							const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
							const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
							return rightTime - leftTime;
						}),
					}
					: null;
			} catch (err) {
				this.errorMessage = getErrorMessage(err, 'Unable to load journal details. Please try again.');
			} finally {
				this.loading = false;
			}
		},
		async deleteJournal(this: any) {
			if (!this.journal || this.journalDeleting) return;
			if (typeof window !== 'undefined' && !window.confirm('Delete this journal? This cannot be undone.')) return;
			this.journalDeleting = true;
			try {
				await journalClient.delete(this.journalId);
				window.location.href = '/journal';
			} catch (err) {
				this.errorMessage = getErrorMessage(err, 'Unable to delete journal.');
			} finally {
				this.journalDeleting = false;
			}
		},
		hasError(this: any) {
			return this.errorMessage !== '';
		},
	};
}

import type { JournalClient } from '../client/journal';

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
				this.journal = envelope.data ?? null;
			} catch (err) {
				this.errorMessage = err instanceof Error ? err.message : 'Unable to load journal details. Please try again.';
			} finally {
				this.loading = false;
			}
		},
		hasError(this: any) {
			return this.errorMessage !== '';
		},
	};
}

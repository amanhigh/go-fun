import { getErrorMessage } from '../../../lib/error';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function normalizeJournal(journal: any) {
	if (!journal) return null;
	return {
		...journal,
		images: journal.images ?? [],
		tags: journal.tags ?? [],
		notes: journal.notes ?? [],
	};
}

export function NewCurrentJournalConcern(pg: JournalDetailPageProvider) {
	return {
		journalId: '',
		journal: null,
		loading: true,
		errorMessage: '',

		hasError(this: any) { return this.errorMessage !== ''; },

		async loadJournal(this: any) {
			this.loading = true;
			this.errorMessage = '';
			try {
				const envelope = await pg().client.get(this.journalId);
				this.journal = normalizeJournal(envelope.data);
				pg().sidebar.tags.sync(this.journal?.tags);
				pg().sidebar.notes.sync(this.journal?.notes);
			} catch (err) {
				this.errorMessage = getErrorMessage(err, 'Unable to load journal details. Please try again.');
			} finally {
				this.loading = false;
			}
		},
	};
}

import { getErrorMessage } from '../../../shared/error';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function normalizeJournal(journal: any) {
	if (!journal) return null;
	return {
		...journal,
		images: journal.images ?? [],
		tags: journal.tags ?? [],
		notes: [...(journal.notes ?? [])].sort((left: any, right: any) => {
			const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
			const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
			return rightTime - leftTime;
		}),
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

import type { JournalDetail } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function normalizeJournal(journal: unknown): JournalDetail | null {
	if (!journal) return null;
	const src = journal as Record<string, unknown>;
	return {
		...src,
		images: (src.images as JournalDetail['images']) ?? [],
		tags: (src.tags as JournalDetail['tags']) ?? [],
		notes: (src.notes as JournalDetail['notes']) ?? [],
	} as JournalDetail;
}

export function NewJournalConcern(pg: JournalDetailPageProvider) {
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
				this.errorMessage = (err as Error).message;
			} finally {
				this.loading = false;
			}
		},
	};
}

import { createLoader } from '../../../lib/loader';
import type { Loader } from '../../../lib/loader';
import type { JournalDetail } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

function normalizeJournal(journal: JournalDetail): JournalDetail {
	return {
		...journal,
		images: journal.images ?? [],
		tags: journal.tags ?? [],
		notes: journal.notes ?? [],
	};
}

export function NewJournalConcern(pg: JournalDetailPageProvider) {
	return {
		journalId: '',
		journal: null,
		loader: createLoader(),

		async loadJournal(this: any) {
			const data = await this.loader.loadData(
				() => pg().client.get(this.journalId),
			);

			if (!data) return;

			this.journal = normalizeJournal(data);
			pg().sidebar.tags.sync(this.journal.tags);
			pg().sidebar.notes.sync(this.journal.notes);
		},
	};
}

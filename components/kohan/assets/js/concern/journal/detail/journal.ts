import { createLoader } from '../../../lib/loader';
import type { Loader } from '../../../lib/loader';
import type { JournalDetail } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

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
		loader: createLoader(true) as Loader,

		async loadJournal(this: any) {
			const data = await this.loader.loadData(
				() => pg().client.get(this.journalId),
			);

			if (!data) return;

			this.journal = normalizeJournal(data);
			pg().sidebar.tags.sync(this.journal?.tags);
			pg().sidebar.notes.sync(this.journal?.notes);
		},
	};
}

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
		detail: null,
		loader: createLoader(),

		async loadJournal(this: any, id: string) {
			await this.loader.load(
				() => pg().client.get(id),
				(data: any) => {
					this.detail = normalizeJournal(data);
					pg().sidebar.tags.sync(this.detail.tags);
					pg().sidebar.notes.sync(this.detail.notes);
				},
			);
		},
	};
}

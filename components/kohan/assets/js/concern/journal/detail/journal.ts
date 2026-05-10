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
			const data = await this.loader.loadData(
				() => pg().client.get(id),
			);

			if (!data) return;

			this.detail = normalizeJournal(data);
			pg().sidebar.tags.sync(this.detail.tags);
			pg().sidebar.notes.sync(this.detail.notes);
		},
	};
}

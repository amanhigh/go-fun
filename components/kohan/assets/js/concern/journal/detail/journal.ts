import { createLoader } from '../../../lib/loader';
import type { Loader } from '../../../lib/loader';
import type { JournalDetail } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

const JOURNAL_OPENED_EVENT = 'kohan:journal-opened';
const JOURNAL_BRIDGE_SELECTOR = '[data-journal-bridge="journal-detail"]';

function publishJournalOpened(journal: JournalDetail): void {
	const bridge = document.querySelector<HTMLElement>(JOURNAL_BRIDGE_SELECTOR);
	if (bridge) bridge.dataset.journalTicker = journal.ticker;
	document.dispatchEvent(new CustomEvent(JOURNAL_OPENED_EVENT, { detail: journal.ticker }));
}

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
					publishJournalOpened(this.detail);
					const index = pg().images.secondSetIndex();
					if (index >= 0) pg().preview.open(index);
				},
			);
		},
	};
}

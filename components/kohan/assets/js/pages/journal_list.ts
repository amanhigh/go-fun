import { NewJournalClient } from '../client/journal';
import type { JournalPageData } from '../types/journal_list_state';
import { createJournalPresentation } from '../concern/journal/common/presentation';
import { createJournalFilterUrlConcern } from '../concern/journal/list/filter_url';
import { createJournalFilter } from '../concern/journal/list/filter';
import { createPresetConcern } from '../concern/journal/list/presets';
import { createPaginationConcern } from '../concern/journal/list/pagination';
import { createJournalTableConcern } from '../concern/journal/list/table';
import '../types/platform';

function createJournalPageData() {
	let page = {} as JournalPageData;
	const pg = () => page;

	page.client = NewJournalClient();
	page.presentation = createJournalPresentation();
	page.table = createJournalTableConcern(pg);
	page.pagination = createPaginationConcern(pg);
	page.presets = createPresetConcern(pg);
	page.filter = createJournalFilter(pg);
	page.filterUrl = createJournalFilterUrlConcern(page);
	page.init = function init(this: any) {
		page = this;
		this.filterUrl.urlToFilter();
		this.presets.syncDatePreset();
		this.presets.syncActiveReviewPreset();
		void this.table.loadJournals();
	};

	return page;
}

document.addEventListener('alpine:init', () => {
	window.Alpine.data('journalPage', createJournalPageData);
});

export {};

import { NewJournalClient } from '../client/journal';
import type { JournalPageData } from '../types/journal_list_state';
import { createJournalPresentation } from '../concern/journal/common/presentation';
import { createJournalFilterUrlConcern } from '../concern/journal/list/filter_url';
import { createJournalFilter } from '../concern/journal/list/filter';
import { createPresetConcern } from '../concern/journal/list/presets';
import { createPaginationState } from '../concern/journal/list/pagination';
import { createJournalTableConcern } from '../concern/journal/list/table';
import '../types/platform';

const journalPageSize = 10;

function createJournalPageData() {
	const client = NewJournalClient();
	const page = {} as JournalPageData;

	page.presentation = createJournalPresentation();
	page.filter = createJournalFilter(page);
	page.filterUrl = createJournalFilterUrlConcern(page.filter);
	page.presets = createPresetConcern(page);
	page.table = createJournalTableConcern(page, client);
	page.pagination = createPaginationState(page, journalPageSize);
	page.init = function init(this: any) {
		console.error('journalPage:init', window.location.search);
		this.filterUrl.urlToFilter();
		this.presets.syncActiveReviewPreset();
		void this.table.loadJournals();
	};

	return page;
}

document.addEventListener('alpine:init', () => {
	window.Alpine.data('journalPage', createJournalPageData);
});

export {};

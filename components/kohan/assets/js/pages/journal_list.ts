import { NewJournalClient } from '../client/journal';
import type { JournalPageData } from '../types/journal_list_concern';
import { newPresentationConcern } from '../concern/journal/common/presentation';
import { newFilterUrlConcern } from '../concern/journal/list/filter_url';
import { newFilterConcern } from '../concern/journal/list/filter';
import { newPresetConcern } from '../concern/journal/list/presets';
import { newPaginationConcern } from '../concern/journal/list/pagination';
import { newTableConcern } from '../concern/journal/list/table';
import '../types/platform';

function createJournalPageData() {
	let page = {} as JournalPageData;
	const pg = () => page;

	page.client = NewJournalClient();
	page.presentation = newPresentationConcern();
	page.table = newTableConcern(pg);
	page.pagination = newPaginationConcern(pg);
	page.presets = newPresetConcern(pg);
	page.filter = newFilterConcern(pg);
	page.filterUrl = newFilterUrlConcern(page);
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

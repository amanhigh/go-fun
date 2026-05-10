import { NewJournalClient } from '../client/journal';
import type { JournalPageData } from '../types/journal/list';
import { NewPresentationConcern } from '../concern/present/factory';
import { NewFilterUrlConcern } from '../concern/journal/list/filter_url';
import { NewFilterConcern } from '../concern/journal/list/filter';
import { NewPresetConcern } from '../concern/journal/list/presets';
import { NewPaginationConcern } from '../concern/journal/list/pagination';
import { NewTableConcern } from '../concern/journal/list/table';
import '../types/platform';

function createJournalPageData() {
	let page = {} as JournalPageData;
	const pg = () => page;

	page.client = NewJournalClient();
	page.present = NewPresentationConcern();
	page.table = NewTableConcern(pg);
	page.pagination = NewPaginationConcern(pg);
	page.presets = NewPresetConcern(pg);
	page.filter = NewFilterConcern(pg);
	page.filterUrl = NewFilterUrlConcern(pg);
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

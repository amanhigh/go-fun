import { NewJournalClient } from '../client/journal';
import type { JournalPageData } from '../types/journal_list_state';
import { createJournalFilter } from '../concern/journal/list/filter';
import { createPresetConcern } from '../concern/journal/list/presets';
import { createJournalListFormatters } from '../concern/journal/list/formatters';
import { createPaginationState } from '../concern/journal/list/pagination';
import { createJournalTableConcern } from '../concern/journal/list/table';
import '../types/platform';

const journalPageSize = 10;

function createJournalPageData() {
	const client = NewJournalClient();
	const page = {} as JournalPageData;

	const formatters = createJournalListFormatters();
	page.table = createJournalTableConcern(page, client);
	page.pagination = createPaginationState(page, journalPageSize);
	page.filter = createJournalFilter(page);
	page.presets = createPresetConcern(page);
	page.init = function init() {
		page.filter.urlToFilter();
		page.presets.syncActiveReviewPreset();
		void page.table.loadJournals();
	};

	page.normalizeStatus = formatters.normalizeStatus;
	page.statusBadgeClass = formatters.statusBadgeClass;
	page.typeBadgeClass = formatters.typeBadgeClass;
	page.formatTimestamp = formatters.formatTimestamp;

	return page;
}

document.addEventListener('alpine:init', () => {
	window.Alpine.data('journalPage', createJournalPageData);
});

export {};

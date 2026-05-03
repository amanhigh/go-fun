import type { JournalClient } from '../../../client/journal';
import type { JournalPageData, JournalTableState } from '../../../types/journal_list_state';

export function createJournalTableConcern(page: JournalPageData, client: JournalClient): JournalTableState {
	const table: JournalTableState = {
		journals: [],
		requestCounter: 0,
		loading: false,
		errorMessage: '',
		applyFilters() {
			page.pagination.resetPage();
			page.filter.filterToUrl();
			void table.loadJournals();
		},
		applyManualFilters() {
			page.presets.clearActiveReviewPreset();
			table.applyFilters();
		},
		async loadJournals() {
			table.loading = true;
			table.errorMessage = '';

			try {
				const response = await client.list(page.pagination.getOffset(), page.pagination.getPageSize(), page.filter.toQueryParams());
				const data = response.data ?? {};
				table.journals = data.journals ?? [];
				page.pagination.setTotalItems(data.metadata?.total ?? table.journals.length);
				page.pagination.setPageFromOffset(data.metadata?.offset ?? 0);
			} finally {
				table.loading = false;
			}
		},
		hasError() { return table.errorMessage !== ''; },
		isEmpty() { return table.journals.length === 0; },
	};

	return table;
}

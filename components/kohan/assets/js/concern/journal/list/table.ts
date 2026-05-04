import type { JournalClient } from '../../../client/journal';
import type { JournalPageData, JournalTableState } from '../../../types/journal_list_state';

export function createJournalTableConcern(page: JournalPageData, client: JournalClient): JournalTableState {
	const table: JournalTableState = {
		journals: [],
		requestCounter: 0,
		loading: false,
		errorMessage: '',
		applyFilters(this: JournalTableState) {
			page.pagination.resetPage();
			page.filterUrl.filterToUrl();
			void this.loadJournals();
		},
		applyManualFilters(this: JournalTableState) {
			page.presets.clearActiveReviewPreset();
			this.applyFilters();
		},
		async loadJournals(this: JournalTableState) {
			this.loading = true;
			this.errorMessage = '';

			try {
				const response = await client.list(page.pagination.getOffset(), page.pagination.getPageSize(), page.filter);
				const data = response.data ?? {};
				this.journals = data.journals ?? [];
				page.pagination.setTotalItems(data.metadata?.total ?? this.journals.length);
				page.pagination.setPageFromOffset(data.metadata?.offset ?? 0);
			} finally {
				this.loading = false;
			}
		},
		hasError(this: JournalTableState) { return this.errorMessage !== ''; },
		isEmpty(this: JournalTableState) { return this.journals.length === 0; },
	};

	return table;
}

import type { JournalPageProvider, JournalTableState } from '../../../types/journal_list_state';

export function createJournalTableConcern(pg: JournalPageProvider): JournalTableState {
	const table: JournalTableState = {
		journals: [],
		requestCounter: 0,
		loading: false,
		errorMessage: '',
		applyFilters(this: JournalTableState) {
			pg().pagination.resetPage();
			pg().filterUrl.filterToUrl();
			void this.loadJournals();
		},
		applyManualFilters(this: JournalTableState) {
			pg().presets.clearActiveReviewPreset();
			pg().filter.datePreset = '';
			this.applyFilters();
		},
		async loadJournals(this: JournalTableState) {
			this.loading = true;
			this.errorMessage = '';

			try {
				const response = await pg().client.list(pg().pagination.getOffset(), pg().pagination.getPageSize(), pg().filter);
				const data = response.data ?? {};
				this.journals = data.journals ?? [];
				pg().pagination.setTotalItems(data.metadata?.total ?? this.journals.length);
				pg().pagination.setPageFromOffset(data.metadata?.offset ?? 0);
			} finally {
				this.loading = false;
			}
		},
		hasError(this: JournalTableState) { return this.errorMessage !== ''; },
		isEmpty(this: JournalTableState) { return this.journals.length === 0; },
	};

	return table;
}

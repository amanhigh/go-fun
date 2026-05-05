import type { JournalPageProvider, JournalTableConcern } from '../../../types/journal_list_concern';

export function newTableConcern(pg: JournalPageProvider): JournalTableConcern {
	const table: JournalTableConcern = {
		journals: [],
		requestCounter: 0,
		loading: false,
		errorMessage: '',
		applyFilters(this: JournalTableConcern) {
			pg().pagination.resetPage();
			pg().filterUrl.filterToUrl();
			void this.loadJournals();
		},
		applyManualFilters(this: JournalTableConcern) {
			pg().presets.clearActiveReviewPreset();
			pg().filter.datePreset = '';
			this.applyFilters();
		},
		async loadJournals(this: JournalTableConcern) {
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
		hasError(this: JournalTableConcern) { return this.errorMessage !== ''; },
		isEmpty(this: JournalTableConcern) { return this.journals.length === 0; },
	};

	return table;
}

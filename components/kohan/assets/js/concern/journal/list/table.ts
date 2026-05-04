import type { JournalClient } from '../../../client/journal';
import type { JournalPageData, JournalTableState } from '../../../types/journal_list_state';

export function createJournalTableConcern(page: JournalPageData, client: JournalClient): JournalTableState {
	const table: JournalTableState = {
		journals: [],
		requestCounter: 0,
		loading: false,
		errorMessage: '',
		applyFilters(this: JournalTableState) {
			const context = (page as any).__runtime ?? page;
			context.pagination.resetPage();
			context.filterUrl.filterToUrl();
			void this.loadJournals();
		},
		applyManualFilters(this: JournalTableState) {
			const context = (page as any).__runtime ?? page;
			context.presets.clearActiveReviewPreset();
			context.filter.datePreset = '';
			this.applyFilters();
		},
		async loadJournals(this: JournalTableState) {
			const context = (page as any).__runtime ?? page;
			this.loading = true;
			this.errorMessage = '';

			try {
				const response = await client.list(context.pagination.getOffset(), context.pagination.getPageSize(), context.filter);
				const data = response.data ?? {};
				this.journals = data.journals ?? [];
				context.pagination.setTotalItems(data.metadata?.total ?? this.journals.length);
				context.pagination.setPageFromOffset(data.metadata?.offset ?? 0);
			} finally {
				this.loading = false;
			}
		},
		hasError(this: JournalTableState) { return this.errorMessage !== ''; },
		isEmpty(this: JournalTableState) { return this.journals.length === 0; },
	};

	return table;
}

import { JournalClient } from './journal_client';
import { journalQueryKeyMap, journalReverseQueryKeyMap, type Journal } from './journal_models';
import { createFilterTracker, createPaginationTracker } from './journal_state';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalPage>): void;
};

function journalPage() {
	const client = new JournalClient();
	const pagination = createPaginationTracker(10);
	const filterTracker = createFilterTracker();
	const trackedFilters = filterTracker as unknown as Record<string, string>;
	const reverseQueryKeyMap = journalReverseQueryKeyMap as Record<string, keyof typeof filterTracker>;
	return {
		journals: [] as Journal[],
		pagination,
		filterTracker,
		loading: false,
		errorMessage: '',
		async loadJournals() {
			this.loading = true;
			this.errorMessage = '';
			try {
				const resp = await client.list(
					this.pagination.getPage() === 1 ? 0 : (this.pagination.getPage() - 1) * this.pagination.getPageSize(),
					this.pagination.getPageSize(),
					this.filterTracker.toQueryParams(),
				);
				const data = resp.data ?? {};
				this.journals = data.journals ?? [];
				this.pagination.setTotalItems(data.metadata?.total ?? this.journals.length);
				this.pagination.setPageFromOffset(data.metadata?.offset ?? 0);
			} catch {
				this.errorMessage = 'Unable to load journals. Please try again.';
			} finally {
				this.loading = false;
			}
		},
		urlToFilter() {
			const params = new URLSearchParams(window.location.search);
			params.forEach((value, key) => {
				const filterKey = reverseQueryKeyMap[key];
				if (filterKey) {
					trackedFilters[filterKey] = value;
				}
			});
		},
		filterToUrl() {
			const params = new URLSearchParams();
			Object.entries(filterTracker.toQueryParams()).forEach(([key, value]) => {
				if (value !== '') params.set(journalQueryKeyMap[key] ?? key, value);
			});
			const nextUrl = params.toString() ? `${window.location.pathname}?${params.toString()}` : window.location.pathname;
			window.history.replaceState({}, '', nextUrl);
		},
		hasError() {
			return this.errorMessage !== '';
		},
		isEmpty() {
			return this.journals.length === 0;
		},
		formatTimestamp(value: string) {
			if (!value) return '—';
			return new Date(value).toLocaleString();
		},
		async prevPage() {
			if (!this.pagination.hasPrev()) return;
			this.pagination.prevPage();
			await this.loadJournals();
		},
		async nextPage() {
			if (!this.pagination.hasNext()) return;
			this.pagination.nextPage();
			await this.loadJournals();
		},
		applyFilters() {
			this.pagination.resetPage();
			this.filterToUrl();
			void this.loadJournals();
		},
		toggleSort(field: string) {
			if (this.filterTracker.sortBy !== field) {
				this.filterTracker.sortBy = field;
				this.filterTracker.sortOrder = 'asc';
				this.applyFilters();
				return;
			}
			this.filterTracker.sortOrder = this.filterTracker.sortOrder === 'asc' ? 'desc' : 'asc';
			this.applyFilters();
		},
		clearFilters() {
			this.filterTracker.clear();
			this.filterToUrl();
			this.applyFilters();
		},
		init() {
			this.urlToFilter();
			void this.loadJournals();
		},
	};
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalPage', journalPage);
});

export {};

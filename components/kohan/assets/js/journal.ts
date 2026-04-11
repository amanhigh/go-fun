import { JournalClient } from './journal_client';
import { journalQueryKeyMap, journalReverseQueryKeyMap, type Journal } from './journal_models';
import { createFilterTracker, createPaginationTracker } from './journal_state';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalPage>): void;
};

const defaultBadgeClass = 'border-slate-300 bg-slate-50 text-slate-700';

const statusBadgeClassMap: Record<string, string> = {
	SUCCESS: 'border-emerald-300 bg-emerald-50 text-emerald-800',
	FAIL: 'border-rose-300 bg-rose-50 text-rose-800',
	RUNNING: 'border-sky-300 bg-sky-50 text-sky-800',
	SET: 'border-amber-300 bg-amber-50 text-amber-800',
	REJECTED: 'border-violet-300 bg-violet-50 text-violet-800',
};

const typeBadgeClassMap: Record<string, string> = {
	REJECTED: 'border-rose-300 bg-rose-50 text-rose-800',
	RESULT: 'border-emerald-300 bg-emerald-50 text-emerald-800',
	SET: 'border-indigo-300 bg-indigo-50 text-indigo-800',
};

const normalizeTag = (value: string): string => (value ?? '').trim().toUpperCase();

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
		requestCounter: 0,
		loading: false,
		errorMessage: '',
		normalizeStatus(value: string) {
			return normalizeTag(value);
		},
		statusBadgeClass(value: string) {
			return statusBadgeClassMap[normalizeTag(value)] ?? defaultBadgeClass;
		},
		typeBadgeClass(value: string) {
			return typeBadgeClassMap[normalizeTag(value)] ?? defaultBadgeClass;
		},
		async loadJournals() {
			const requestId = this.requestCounter + 1;
			this.requestCounter = requestId;
			this.loading = true;
			this.errorMessage = '';
			try {
				const resp = await client.list(
					this.pagination.getPage() === 1 ? 0 : (this.pagination.getPage() - 1) * this.pagination.getPageSize(),
					this.pagination.getPageSize(),
					this.filterTracker.toQueryParams(),
				);
				if (requestId !== this.requestCounter) return;
				const data = resp.data ?? {};
				this.journals = data.journals ?? [];
				this.pagination.setTotalItems(data.metadata?.total ?? this.journals.length);
				this.pagination.setPageFromOffset(data.metadata?.offset ?? 0);
			} catch {
				if (requestId !== this.requestCounter) return;
				this.errorMessage = 'Unable to load journals. Please try again.';
			} finally {
				if (requestId !== this.requestCounter) return;
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
			const parsed = new Date(value);
			if (Number.isNaN(parsed.getTime())) return '—';
			return parsed.toLocaleString();
		},
		toDateInputValue(date: Date) {
			const year = date.getFullYear();
			const month = `${date.getMonth() + 1}`.padStart(2, '0');
			const day = `${date.getDate()}`.padStart(2, '0');
			return `${year}-${month}-${day}`;
		},
		applyCreatedPreset(preset: string) {
			const today = new Date();
			const endDate = this.toDateInputValue(today);
			const daysMap: Record<string, number> = { today: 0, last7: 7, last30: 30 };
			const days = daysMap[preset] ?? 7;
			if (days === 0) {
				this.filterTracker.createdAfter = endDate;
				this.filterTracker.createdBefore = endDate;
				this.applyFilters();
				return;
			}
			const startDate = new Date(today);
			startDate.setDate(today.getDate() - days);
			this.filterTracker.createdAfter = this.toDateInputValue(startDate);
			this.filterTracker.createdBefore = endDate;
			this.applyFilters();
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
		onCreatedDateChange() {
			this.filterTracker.createdBefore = this.filterTracker.createdAfter;
			this.applyFilters();
		},
		toggleSort(field: string) {
			this.filterTracker.sortOrder = this.filterTracker.sortBy !== field
				? 'asc'
				: this.filterTracker.sortOrder === 'asc' ? 'desc' : 'asc';
			this.filterTracker.sortBy = field;
			this.applyFilters();
		},
		clearFilters() {
			this.filterTracker.clear();
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

import { JournalClient } from './journal_client';
import type { Journal } from './journal_models';
import { createJournalListFilterActions, createReviewPresets } from './journal_list_filters';
import { createJournalListFormatters } from './journal_list_formatters';
import { createFilterTracker, createPaginationTracker } from './journal_state';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalPage>): void;
};

function journalPage() {
	const client = new JournalClient();
	const pagination = createPaginationTracker(10);
	const filterTracker = createFilterTracker();
	const reviewPresets = createReviewPresets();

	return {
		journals: [] as Journal[],
		reviewPresets,
		activeReviewPreset: '',
		pagination,
		filterTracker,
		requestCounter: 0,
		loading: false,
		errorMessage: '',
		...createJournalListFormatters(),
		...createJournalListFilterActions(filterTracker as unknown as Record<string, string>),
		async loadJournals(this: any) {
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
		hasError(this: any) {
			return this.errorMessage !== '';
		},
		isEmpty(this: any) {
			return this.journals.length === 0;
		},
		async prevPage(this: any) {
			if (!this.pagination.hasPrev()) return;
			this.pagination.prevPage();
			await this.loadJournals();
		},
		async nextPage(this: any) {
			if (!this.pagination.hasNext()) return;
			this.pagination.nextPage();
			await this.loadJournals();
		},
		init(this: any) {
			this.urlToFilter();
			this.syncActiveReviewPreset();
			void this.loadJournals();
		},
	};
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalPage', journalPage);
});

export {};

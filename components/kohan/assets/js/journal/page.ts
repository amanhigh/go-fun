import { NewJournalClient, type Journal } from '../client/journal';
import { createFilterActions } from './filter_actions';
import { createFilterPresetActions, createReviewPresets } from './filter_presets';
import { createFilterUrlActions } from './filter_url';
import { createJournalListFormatters } from './formatters';
import { createJournalFilterState } from './filter_state';
import { createPaginationState } from './pagination';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalPage>): void;
};

function journalPage() {
	const client = NewJournalClient();
	const pagination = createPaginationState(10);
	const filter = createJournalFilterState();
	const reviewPresets = createReviewPresets();

	return {
		journals: [] as Journal[],
		reviewPresets,
		activeReviewPreset: '',
		pagination,
		filter,
		requestCounter: 0,
		loading: false,
		errorMessage: '',
		...createJournalListFormatters(),
		...createFilterPresetActions(),
		...createFilterUrlActions(filter),
		...createFilterActions(filter),
		async loadJournals(this: any) {
			const requestId = this.requestCounter + 1;
			this.requestCounter = requestId;
			this.loading = true;
			this.errorMessage = '';
			try {
			const resp = await client.list(
					this.pagination.getPage() === 1 ? 0 : (this.pagination.getPage() - 1) * this.pagination.getPageSize(),
					this.pagination.getPageSize(),
					this.filter.toQueryParams(),
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

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
		beginJournalLoad(this: any) {
			this.requestCounter += 1;
			this.loading = true;
			this.errorMessage = '';
			return this.requestCounter;
		},
		isCurrentJournalLoad(this: any, requestId: number) {
			return requestId === this.requestCounter;
		},
		applyJournalResponse(this: any, resp: Awaited<ReturnType<typeof client.list>>) {
			const data = resp.data ?? {};
			this.journals = data.journals ?? [];
			this.pagination.setTotalItems(data.metadata?.total ?? this.journals.length);
			this.pagination.setPageFromOffset(data.metadata?.offset ?? 0);
		},
		async loadJournals(this: any) {
			const requestId = this.beginJournalLoad();
			try {
				const resp = await client.list(
					this.pagination.getOffset(),
					this.pagination.getPageSize(),
					this.filter.toQueryParams(),
				);
				if (!this.isCurrentJournalLoad(requestId)) return;
				this.applyJournalResponse(resp);
			} catch {
				if (!this.isCurrentJournalLoad(requestId)) return;
				this.errorMessage = 'Unable to load journals. Please try again.';
			} finally {
				if (!this.isCurrentJournalLoad(requestId)) return;
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

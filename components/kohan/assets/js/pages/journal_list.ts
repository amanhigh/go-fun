import { NewJournalClient } from '../client/journal';
import type { CreateJournalPageStateInput, JournalPageState } from '../types/journal_list_state';
import type { JournalClient } from '../client/journal';
import { createFilterActions, createJournalFilter } from '../concern/journal/list/filter';
import { buildReviewPresetList, createPresetActions } from '../concern/journal/list/presets';
import { createJournalListFormatters } from '../concern/journal/list/formatters';
import { createPaginationState } from '../concern/journal/list/pagination';

declare const Alpine: {
	data(name: string, callback: () => any): void;
};

function createJournalPageState(input: CreateJournalPageStateInput): JournalPageState {
	return {
		journals: [],
		reviewPresets: input.reviewPresets,
		activeReviewPreset: '',
		pagination: input.pagination,
		filter: input.filter,
		requestCounter: 0,
		loading: false,
		errorMessage: '',
	};
}

function createJournalPageActions(client: JournalClient) {
	async function loadJournals(this: any) {
		this.loading = true;
		this.errorMessage = '';

		try {
			const response = await client.list(this.pagination.getOffset(), this.pagination.getPageSize(), this.filter.toQueryParams());
			const data = response.data ?? {};
			this.journals = data.journals ?? [];
			this.pagination.setTotalItems(data.metadata?.total ?? this.journals.length);
			this.pagination.setPageFromOffset(data.metadata?.offset ?? 0);
		} finally {
			this.loading = false;
		}
	}

	function applyFilters(this: any) {
		this.pagination.resetPage();
		this.filterToUrl();
		void this.loadJournals();
	}

	function applyManualFilters(this: any) {
		this.clearActiveReviewPreset();
		this.applyFilters();
	}

	return {
		applyFilters,
		applyManualFilters,
		loadJournals,
		hasError(this: any) { return this.errorMessage !== ''; },
		isEmpty(this: any) { return this.journals.length === 0; },
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

function journalPage() {
	const client = NewJournalClient();
	const pagination = createPaginationState(10);
	const filter = createJournalFilter();
	const reviewPresets = buildReviewPresetList();
	const formatters = createJournalListFormatters();
	const filterActions = createFilterActions();
	const state = createJournalPageState({ filter, pagination, reviewPresets });
	const presetActions = createPresetActions();
	const pageActions = createJournalPageActions(client);

	return Object.assign(state, formatters, presetActions, pageActions, filterActions);
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalPage', journalPage);
});

export {};

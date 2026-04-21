import { NewJournalClient, type Journal } from '../client/journal';
import { createFilterActions } from './filter_actions';
import { createReviewPresets } from './filter_presets';
import { createFilterUrlActions } from './filter_url';
import { createJournalListFormatters } from './formatters';
import { createJournalFilterState } from './filter_state';
import { createJournalPageActions } from './page_actions';
import { createPaginationState } from './pagination';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalPage>): void;
};

function journalPage() {
	const client = NewJournalClient();
	const pagination = createPaginationState(10);
	const filter = createJournalFilterState();
	const reviewPresets = createReviewPresets();
	const urlActions = createFilterUrlActions(filter);
	const state = {
		journals: [] as Journal[],
		reviewPresets,
		activeReviewPreset: '',
		pagination,
		filter,
		requestCounter: 0,
		loading: false,
		errorMessage: '',
	};
	const pageActions = createJournalPageActions({
		client,
		filter,
		pagination,
		reviewPresets,
		state,
		filterToUrl: urlActions.filterToUrl,
		urlToFilter: urlActions.urlToFilter,
	});

	return Object.assign(
		state,
		createJournalListFormatters(filter),
		urlActions,
		pageActions,
		createFilterActions({ filter, applyManualFilters: pageActions.applyManualFilters }),
	);
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalPage', journalPage);
});

export {};

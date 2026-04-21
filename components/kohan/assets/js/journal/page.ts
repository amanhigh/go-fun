import { NewJournalClient } from '../client/journal';
import { createFilterActions } from './filter_actions';
import { buildReviewPresetList } from './presets';
import { createJournalListFormatters } from './formatters';
import { createJournalFilter } from './filter';
import { createJournalPageActions } from './page_actions';
import { createJournalPageState } from './page_state';
import { createPresetActions } from './preset_actions';
import { createPaginationState } from './pagination';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalPage>): void;
};

function journalPage() {
	const client = NewJournalClient();
	const pagination = createPaginationState(10);
	const filter = createJournalFilter();
	const reviewPresets = buildReviewPresetList();
	const formatters = createJournalListFormatters(filter);
	const filterActions = createFilterActions({ filter, applyManualFilters: () => pageActions.applyManualFilters() });
	const state = createJournalPageState({ filter, pagination, reviewPresets });
	const presetActions = createPresetActions({
		filter,
		state,
		applyFilters: () => pageActions.applyFilters(),
		clearFilters: () => filter.clear(),
	});
	const pageActions = createJournalPageActions({
		client,
		filter,
		pagination,
		state,
		filterToUrl: filterActions.filterToUrl,
		urlToFilter: filterActions.urlToFilter,
		clearActiveReviewPreset: presetActions.clearActiveReviewPreset,
		syncActiveReviewPreset: presetActions.syncActiveReviewPreset,
	});

	return Object.assign(state, formatters, presetActions, pageActions, filterActions);
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalPage', journalPage);
});

export {};

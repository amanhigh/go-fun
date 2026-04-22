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
	data(name: string, callback: () => any): void;
};

function journalPage() {
	const client = NewJournalClient();
	const pagination = createPaginationState(10);
	const filter = createJournalFilter();
	const reviewPresets = buildReviewPresetList();
	const formatters = createJournalListFormatters();
	const filterActions = createFilterActions();
	const state = createJournalPageState({ filter, pagination, reviewPresets });
	const presetActions = createPresetActions();
	const pageActions = createJournalPageActions({ client });

	return Object.assign(state, formatters, presetActions, pageActions, filterActions);
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalPage', journalPage);
});

export {};

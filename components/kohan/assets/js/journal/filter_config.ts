import type { JournalFilterKey } from './filter_state';

export const journalQueryMap: Record<JournalFilterKey, string> = {
	ticker: 'search',
	createdAfter: 'created-after',
	createdBefore: 'created-before',
	sortBy: 'sort-by',
	sortOrder: 'sort-order',
};

export const journalReverseMap: Record<string, JournalFilterKey> = {
	search: 'ticker',
	ticker: 'ticker',
	type: 'type',
	status: 'status',
	sequence: 'sequence',
	'created-after': 'createdAfter',
	'created-before': 'createdBefore',
	reviewed: 'reviewed',
	'sort-by': 'sortBy',
	'sort-order': 'sortOrder',
};

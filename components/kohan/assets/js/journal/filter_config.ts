export const journalQueryKeyMap: Record<string, string> = {
	ticker: 'search',
	createdAfter: 'created-after',
	createdBefore: 'created-before',
	sortBy: 'sort-by',
	sortOrder: 'sort-order',
};

export const journalReverseQueryKeyMap: Record<string, string> = {
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

export type Journal = {
	id: string;
	ticker: string;
	sequence: string;
	type: string;
	status: string;
	created_at: string;
	reviewed_at?: string | null;
};

export type JournalList = {
	journals?: Journal[];
	metadata?: {
		total?: number;
		offset?: number;
		limit?: number;
	};
};

export type JournalListFilters = {
	ticker?: string;
	type?: string;
	status?: string;
	sequence?: string;
	createdAfter?: string;
	createdBefore?: string;
	reviewed?: string;
	sortBy?: string;
	sortOrder?: string;
};

export type Envelope<T> = {
	status: string;
	data: T;
};

export const journalQueryKeyMap: Record<string, string> = {
	createdAfter: 'created-after',
	createdBefore: 'created-before',
	sortBy: 'sort-by',
	sortOrder: 'sort-order',
};

export const journalReverseQueryKeyMap: Record<string, string> = {
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

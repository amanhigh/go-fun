export type Journal = {
	id: string;
	ticker: string;
	sequence: string;
	type: string;
	status: string;
	created_at: string;
	reviewed_at?: string | null;
	images?: JournalImage[];
	tags?: JournalTag[];
	notes?: JournalNote[];
	deleted_at?: string | null;
};

export type JournalImage = {
	id: string;
	timeframe: string;
	file_name: string;
	created_at?: string;
};

export type JournalTag = {
	id: string;
	tag: string;
	type?: string;
	override?: string;
	created_at?: string;
};

export type JournalNote = {
	id: string;
	status: string;
	content: string;
	format?: string;
	created_at?: string;
};

export type JournalReviewUpdate = {
	status?: string;
	reviewed_at: string | null;
};

export type JournalReviewStatusResponse = {
	id: string;
	status: string;
	reviewed_at: string | null;
};

export type JournalNoteCreate = {
	status: string;
	content: string;
	format: 'MARKDOWN' | 'PLAINTEXT';
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

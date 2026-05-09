// ===== Request Types =====

export type JournalFilterKey = 'ticker' | 'type' | 'status' | 'sequence' | 'createdAfter' | 'createdBefore' | 'reviewed' | 'sortBy' | 'sortOrder';

export type JournalFilters = Record<JournalFilterKey, string>;

export type JournalListRequest = Partial<JournalFilters>;

export type JournalUpdateRequest = {
	status?: string;
	reviewed_at: string | null;
};

export type JournalNoteRequest = {
	status: string;
	content: string;
	format: 'MARKDOWN' | 'PLAINTEXT';
};

export type JournalTagRequest = {
	tag: string;
	type: string;
	override?: string;
};

// ===== Response Types =====

export type JournalTimeframe = 'DL' | 'WK' | 'MN' | 'TMN' | 'SMN' | 'YR';

export type JournalImage = {
	id: string;
	timeframe: JournalTimeframe;
	file_name: string;
	created_at?: string;
};

export type JournalNote = {
	id: string;
	status: string;
	content: string;
	format?: string;
	created_at?: string;
};

export type JournalTag = {
	id: string;
	tag: string;
	type?: string;
	override?: string;
	created_at?: string;
};

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

export type JournalList = {
	journals?: Journal[];
	metadata?: {
		total?: number;
		offset?: number;
		limit?: number;
	};
};

export type JournalUpdate = {
	id: string;
	status: string;
	reviewed_at: string | null;
};

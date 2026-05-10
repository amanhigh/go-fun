import type { PaginatedResponse } from '../client/base';

// ===== Enums =====

export const JournalType = {
	TAKEN: 'TAKEN',
	REJECTED: 'REJECTED',
} as const;
export type JournalType = (typeof JournalType)[keyof typeof JournalType];

export const JournalStatus = {
	SET: 'SET',
	RUNNING: 'RUNNING',
	SUCCESS: 'SUCCESS',
	FAIL: 'FAIL',
	MISSED: 'MISSED',
	JUST_LOSS: 'JUST_LOSS',
	BROKEN: 'BROKEN',
} as const;
export type JournalStatus = (typeof JournalStatus)[keyof typeof JournalStatus];

export const JournalSequence = {
	MWD: 'MWD',
	YR: 'YR',
	WDH: 'WDH',
} as const;
export type JournalSequence = (typeof JournalSequence)[keyof typeof JournalSequence];

export const JournalTimeframe = {
	DL: 'DL',
	WK: 'WK',
	MN: 'MN',
	TMN: 'TMN',
	SMN: 'SMN',
	YR: 'YR',
} as const;
export type JournalTimeframe = (typeof JournalTimeframe)[keyof typeof JournalTimeframe];

export const JournalTagType = {
	REASON: 'REASON',
	MANAGEMENT: 'MANAGEMENT',
	DIRECTION: 'DIRECTION',
} as const;
export type JournalTagType = (typeof JournalTagType)[keyof typeof JournalTagType];

export const JournalNoteFormat = {
	MARKDOWN: 'MARKDOWN',
	PLAINTEXT: 'PLAINTEXT',
} as const;
export type JournalNoteFormat = (typeof JournalNoteFormat)[keyof typeof JournalNoteFormat];

export const JournalSortBy = {
	CREATED_AT: 'created_at',
	TICKER: 'ticker',
	SEQUENCE: 'sequence',
} as const;
export type JournalSortBy = (typeof JournalSortBy)[keyof typeof JournalSortBy];

export const JournalSortOrder = {
	ASC: 'asc',
	DESC: 'desc',
} as const;
export type JournalSortOrder = (typeof JournalSortOrder)[keyof typeof JournalSortOrder];

// ===== Frontend-owned Consts =====

export const ReviewedFilter = {
	ALL: '' as const,
	PENDING: 'false' as const,
	REVIEWED: 'true' as const,
} as const;
export type ReviewedFilter = (typeof ReviewedFilter)[keyof typeof ReviewedFilter];

// ===== Request Types =====

export type JournalFilterKey = 'ticker' | 'type' | 'status' | 'sequence' | 'createdAfter' | 'createdBefore' | 'reviewed' | 'sortBy' | 'sortOrder';

export type JournalFilters = Record<JournalFilterKey, string>;

export type JournalListRequest = Partial<JournalFilters>;

export type JournalUpdateRequest = {
	status?: JournalStatus;
	reviewed_at?: string | null;
};

export type JournalNoteRequest = {
	status: JournalStatus;
	content: string;
	format: JournalNoteFormat;
};

export type JournalTagRequest = {
	tag: string;
	type: JournalTagType;
	override?: string;
};

// ===== Response Types =====

export type JournalImage = {
	id: string;
	timeframe: JournalTimeframe;
	file_name: string;
	created_at: string;
};

export type JournalNote = {
	id: string;
	status: JournalStatus;
	content: string;
	format: JournalNoteFormat;
	created_at: string;
};

export type JournalTag = {
	id: string;
	tag: string;
	type: JournalTagType;
	override?: string;
	created_at: string;
};

export type Journal = {
	id: string;
	ticker: string;
	sequence: JournalSequence;
	type: JournalType;
	status: JournalStatus;
	created_at: string;
	reviewed_at?: string | null;
	images?: JournalImage[];
	tags?: JournalTag[];
	notes?: JournalNote[];
	deleted_at?: string | null;
};

export type JournalList = {
	journals: Journal[];
	metadata: PaginatedResponse;
};

export type JournalUpdate = {
	id: string;
	status: JournalStatus;
	reviewed_at?: string | null;
};

/** Normalized detail view with always-present association arrays. */
export type JournalDetail = Journal & {
	images: JournalImage[];
	tags: JournalTag[];
	notes: JournalNote[];
};

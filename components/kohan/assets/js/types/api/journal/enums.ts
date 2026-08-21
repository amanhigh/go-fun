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

/** Read-domain journal top timeframe. MN is migration-only. */
export const JournalTopTimeframe = {
	YR: 'YR',
	SMN: 'SMN',
	TMN: 'TMN',
	MN: 'MN',
} as const;
export type JournalTopTimeframe = (typeof JournalTopTimeframe)[keyof typeof JournalTopTimeframe];

/** Filter/URL-exposed top timeframes; excludes migration-only MN. */
export const JournalTopTimeframeFilter = {
	YR: 'YR',
	SMN: 'SMN',
	TMN: 'TMN',
} as const;
export type JournalTopTimeframeFilter = (typeof JournalTopTimeframeFilter)[keyof typeof JournalTopTimeframeFilter];

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
	TOP_TIMEFRAME: 'top_timeframe',
} as const;
export type JournalSortBy = (typeof JournalSortBy)[keyof typeof JournalSortBy];

export const JournalSortOrder = {
	ASC: 'asc',
	DESC: 'desc',
} as const;
export type JournalSortOrder = (typeof JournalSortOrder)[keyof typeof JournalSortOrder];

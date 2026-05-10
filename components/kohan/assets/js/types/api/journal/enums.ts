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

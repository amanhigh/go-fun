import type { JournalStatus, JournalNoteFormat, JournalTagType } from './enums';

// ===== Frontend-owned Query/Filter Consts =====

export const ReviewedFilter = {
	ALL: '' as const,
	PENDING: 'false' as const,
	REVIEWED: 'true' as const,
} as const;
export type ReviewedFilter = (typeof ReviewedFilter)[keyof typeof ReviewedFilter];

export type JournalFilterKey = 'ticker' | 'type' | 'status' | 'sequence' | 'createdAfter' | 'createdBefore' | 'reviewed' | 'sortBy' | 'sortOrder';

export type JournalFilters = Record<JournalFilterKey, string>;

export type JournalListRequest = Partial<JournalFilters>;

// ===== Mutation Request DTOs =====

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

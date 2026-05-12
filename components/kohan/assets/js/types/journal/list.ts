import type { Journal } from '../api/journal/response';
import type { ReviewedFilter } from '../api/journal/request';
import type { JournalType, JournalStatus, JournalSequence, JournalSortBy, JournalSortOrder } from '../api/journal/enums';
import type { Loader } from '../../lib/loader';
import type { Collection } from '../core/collection';
import type { JournalPageBase, PageProvider } from './page';
import type { QuickConcern } from '../core/quick';

// ===== Main Page Composition =====

export type JournalPage = JournalPageBase & {
	filter: JournalFilterConcern;
	filterUrl: JournalFilterUrlConcern;
	pagination: PaginationConcern;
	presets: PresetConcern;
	quick: QuickConcern;
	table: JournalTableConcern;
};

export type JournalPageProvider = PageProvider<JournalPage>;

// ===== Page Sub-Concerns =====

export type PaginationConcern = {
	page: number;
	pageSize: number;
	totalItems: number;
	getPage(): number;
	getPageSize(): number;
	getOffset(): number;
	getTotalItems(): number;
	getTotalPages(): number;
	hasNext(): boolean;
	hasPrev(): boolean;
	setTotalItems(count: number): void;
	setPageFromOffset(offset: number): void;
	nextPage(): void;
	prevPage(): void;
	resetPage(): void;
	previousPage(): Promise<void>;
	nextJournalPage(): Promise<void>;
	summary(): string;
};

export const DatePresetName = {
	ALL: '' as const,
	TODAY: 'today' as const,
	LAST7: 'last7' as const,
	LAST30: 'last30' as const,
} as const;
export type DatePresetName = (typeof DatePresetName)[keyof typeof DatePresetName];
export type NonEmptyDatePresetName = Exclude<DatePresetName, ''>;

export type JournalFilterValues = {
	ticker: string;
	type: JournalType | '';
	status: JournalStatus | '';
	sequence: JournalSequence | '';
	createdAfter: string;
	createdBefore: string;
	reviewed: ReviewedFilter;
	sortBy: JournalSortBy;
	sortOrder: JournalSortOrder;
};

export type JournalFilterConcern = JournalFilterValues & {
	datePreset: DatePresetName;
	clear(): void;
	hasActiveState(): boolean;
	toggleSort(field: JournalSortBy): void;
	applyFilters(): void;
	applyManualFilters(): void;
	clearFilters(): void;
	applyCreatedDate(createdAt: string): void;
};

export type JournalFilterUrlConcern = {
	urlToFilter(): void;
	filterToUrl(): void;
};

export type ReviewPreset = {
	isAnchor: boolean;
	label: string;
	createdAfter: string;
	createdBefore: string;
};

export type PresetConcern = {
	reviewPresets: ReviewPreset[];
	activeReviewPreset: string;
	clearActiveReviewPreset(): void;
	syncActiveReviewPreset(): void;
	syncDatePreset(): void;
	reviewPresetClass(reviewPreset: ReviewPreset): string;
	applyCreatedPreset(preset: NonEmptyDatePresetName): void;
	applyReviewPreset(reviewPreset: ReviewPreset): void;
};

export type JournalTableConcern = Collection<Journal> & {
	loader: Loader;
	load(): Promise<void>;
};

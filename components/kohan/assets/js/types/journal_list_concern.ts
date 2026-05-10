import type { Journal, JournalFilterKey, JournalType, JournalStatus, JournalSequence, JournalSortBy, JournalSortOrder, ReviewedFilter } from './journal_api';
import type { JournalClient } from '../client/journal';
import type { PresentationConcern } from './presentation_concern';
import type { PresentationConcern as PresentConcern } from './present';
import type { Loader } from '../lib/loader';

export type JournalPageProvider = () => JournalPageData;

export type JournalPageData = {
	client: JournalClient;
	presentation: PresentationConcern;
	present: PresentConcern;
	filter: JournalFilterConcern;
	filterUrl: JournalFilterUrlConcern;
	pagination: PaginationConcern;
	presets: PresetConcern;
	table: JournalTableConcern;
	init(): void;
};

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
	toggleType(): void;
	typeToggle(): { label: string; className: string; nextType: JournalType | '' };
	toggleSort(field: JournalSortBy): void;
	applyFilters(): void;
	applyManualFilters(): void;
	clearFilters(): void;
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

export type JournalTableConcern = {
	journals: Journal[];
	loader: Loader;
	loadJournals(): Promise<void>;
	isEmpty(): boolean;
};

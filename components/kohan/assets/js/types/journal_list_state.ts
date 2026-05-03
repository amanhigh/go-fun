import type { Journal, JournalFilterKey } from './journal_api';

export type JournalPageData = JournalListFormatters & {
	filter: JournalFilterState;
	pagination: PaginationState;
	presets: PresetState;
	table: JournalTableState;
	init(): void;
};

export type PaginationState = {
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

export type JournalFilterState = Record<JournalFilterKey, string> & {
	clear(): void;
	toQueryParams(): Record<JournalFilterKey, string>;
	hasActiveState(): boolean;
	urlToFilter(): void;
	filterToUrl(): void;
	toggleType(): void;
	typeToggleLabel(): string;
	typeToggleClass(): string;
	onCreatedDateChange(): void;
	toggleSort(field: string): void;
	applyManualFilters(): void;
	clearFilters(): void;
};

export type ReviewPreset = {
	isAnchor: boolean;
	label: string;
	createdAfter: string;
	createdBefore: string;
};

export type PresetState = {
	reviewPresets: ReviewPreset[];
	activeReviewPreset: string;
	clearActiveReviewPreset(): void;
	syncActiveReviewPreset(): void;
	reviewPresetClass(reviewPreset: ReviewPreset): string;
	applyCreatedPreset(preset: 'today' | 'last7' | 'last30'): void;
	applyReviewPreset(reviewPreset: ReviewPreset): void;
};

export type JournalTableState = {
	journals: Journal[];
	requestCounter: number;
	loading: boolean;
	errorMessage: string;
	applyFilters(): void;
	applyManualFilters(): void;
	loadJournals(): Promise<void>;
	hasError(): boolean;
	isEmpty(): boolean;
};

export type JournalListFormatters = {
	normalizeStatus(value: string): string;
	statusBadgeClass(value: string): string;
	typeBadgeClass(value: string): string;
	formatTimestamp(value: string | null | undefined): string;
};

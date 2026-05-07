import type { Journal, JournalFilterKey } from './journal_api';
import type { JournalClient } from '../client/journal';
import type { PresentationConcern } from './presentation_concern';

export type JournalPageProvider = () => JournalPageData;

export type JournalPageData = {
	client: JournalClient;
	presentation: PresentationConcern;
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

export type DatePresetName = '' | 'today' | 'last7' | 'last30';
export type NonEmptyDatePresetName = Exclude<DatePresetName, ''>;

export type JournalFilterConcern = Record<JournalFilterKey, string> & {
	datePreset: DatePresetName;
	clear(): void;
	hasActiveState(): boolean;
	toggleType(): void;
	typeToggle(): { label: string; className: string; nextType: string };
	toggleSort(field: 'ticker' | 'sequence' | 'created_at'): void;
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
	loading: boolean;
	loadJournals(): Promise<void>;
	isEmpty(): boolean;
};

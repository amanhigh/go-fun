import type { Journal } from './journal';

export type JournalFilterKey = 'ticker' | 'type' | 'status' | 'sequence' | 'createdAfter' | 'createdBefore' | 'reviewed' | 'sortBy' | 'sortOrder';

export type JournalFilters = Record<JournalFilterKey, string>;

export type JournalListRequest = Partial<JournalFilters>;

type FilterConfigEntry = {
	queryKey?: string;
	aliases?: readonly string[];
};

const journalFilterConfig: Record<JournalFilterKey, FilterConfigEntry> = {
	ticker: { queryKey: 'search', aliases: ['ticker'] },
	type: {},
	status: {},
	sequence: {},
	createdAfter: { queryKey: 'created-after' },
	createdBefore: { queryKey: 'created-before' },
	reviewed: {},
	sortBy: { queryKey: 'sort-by' },
	sortOrder: { queryKey: 'sort-order' },
};

export const journalFields: JournalFilterKey[] = ['ticker', 'type', 'status', 'sequence', 'createdAfter', 'createdBefore', 'reviewed', 'sortBy', 'sortOrder'];

export const journalQueryMap: Partial<Record<JournalFilterKey, string>> = journalFields.reduce((queryMap, field) => {
	const entry = journalFilterConfig[field];
	if (!entry.queryKey) return queryMap;
	return { ...queryMap, [field]: entry.queryKey };
}, {} as Partial<Record<JournalFilterKey, string>>);

export const journalReverseMap: Record<string, JournalFilterKey> = journalFields.reduce((reverseMap, field) => {
	const queryKey = journalQueryMap[field] ?? field;
	const aliases = journalFilterConfig[field].aliases ?? [];

	return {
		...reverseMap,
		[queryKey]: field,
		...aliases.reduce<Record<string, JournalFilterKey>>((aliasMap, alias) => ({ ...aliasMap, [alias]: field }), {}),
	};
}, {} as Record<string, JournalFilterKey>);

export const journalFilterUrlMapping = {
	fields: journalFields,
	queryMap: journalQueryMap,
	reverseMap: journalReverseMap,
} as const;

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
};

export type JournalFilterState = Record<JournalFilterKey, string> & {
	clear(): void;
	toQueryParams(): Record<JournalFilterKey, string>;
	hasActiveState(): boolean;
};

export type ReviewPreset = {
	isAnchor: boolean;
	label: string;
	createdAfter: string;
	createdBefore: string;
};

export type JournalPageState = {
	journals: Journal[];
	reviewPresets: ReviewPreset[];
	activeReviewPreset: string;
	pagination: PaginationState;
	filter: JournalFilterState;
	requestCounter: number;
	loading: boolean;
	errorMessage: string;
};

export type CreateJournalPageStateInput = {
	filter: JournalFilterState;
	pagination: PaginationState;
	reviewPresets: ReviewPreset[];
};

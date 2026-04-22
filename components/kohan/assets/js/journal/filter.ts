interface FilterConfigEntry {
	queryKey?: string;
	aliases?: readonly string[];
}

const journalFilterConfig: Record<string, FilterConfigEntry> = {
	ticker: {
		queryKey: 'search',
		aliases: ['ticker'],
	},
	type: {},
	status: {},
	sequence: {},
	createdAfter: {
		queryKey: 'created-after',
	},
	createdBefore: {
		queryKey: 'created-before',
	},
	reviewed: {},
	sortBy: {
		queryKey: 'sort-by',
	},
	sortOrder: {
		queryKey: 'sort-order',
	},
};

export type JournalFilterKey = 'ticker' | 'type' | 'status' | 'sequence' | 'createdAfter' | 'createdBefore' | 'reviewed' | 'sortBy' | 'sortOrder';

export type JournalFilters = Record<JournalFilterKey, string>;

export interface JournalFilterState extends JournalFilters {
	clear(): void;
	toQueryParams(): JournalFilters;
	hasActiveState(): boolean;
}

export const journalFields: JournalFilterKey[] = ['ticker', 'type', 'status', 'sequence', 'createdAfter', 'createdBefore', 'reviewed', 'sortBy', 'sortOrder'];

function createDefaultJournalFilters(): JournalFilters {
	return journalFields.reduce<JournalFilters>((defaults, field) => ({
		...defaults,
		[field]: '',
	}), {} as JournalFilters);
}

export const journalQueryMap: Partial<Record<JournalFilterKey, string>> = journalFields.reduce((queryMap, field) => {
	const entry = journalFilterConfig[field];
	if (!entry.queryKey) {
		return queryMap;
	}

	return {
		...queryMap,
		[field]: entry.queryKey,
	};
}, {} as Partial<Record<JournalFilterKey, string>>);

export const journalReverseMap: Record<string, JournalFilterKey> = journalFields.reduce((reverseMap, field) => {
	const queryKey = journalQueryMap[field] ?? field;
	const aliases = journalFilterConfig[field].aliases ?? [];

	return {
		...reverseMap,
		[queryKey]: field,
		...aliases.reduce<Record<string, JournalFilterKey>>((aliasMap, alias) => ({
			...aliasMap,
			[alias]: field,
		}), {}),
	};
}, {} as Record<string, JournalFilterKey>);

export const journalFilterUrlMapping = {
	fields: journalFields,
	queryMap: journalQueryMap,
	reverseMap: journalReverseMap,
} as const;

export function createJournalFilter(): JournalFilterState {
	const state = { ...createDefaultJournalFilters() } as JournalFilterState;

	state.clear = function clear() {
		Object.assign(state, createDefaultJournalFilters());
	};

	state.toQueryParams = function toQueryParams() {
		return { ...state };
	};

	state.hasActiveState = function hasActiveState() {
		return journalFields.some((field) => state[field] !== '');
	};

	return state;
}

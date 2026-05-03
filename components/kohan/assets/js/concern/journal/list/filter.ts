import { syncStateToUrl, syncUrlToState } from '../../../shared/url_state';
import type { JournalFilterKey, JournalFilters } from '../../../types/journal_api';
import type { JournalFilterState, JournalPageData } from '../../../types/journal_list_state';

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

const journalFields: JournalFilterKey[] = ['ticker', 'type', 'status', 'sequence', 'createdAfter', 'createdBefore', 'reviewed', 'sortBy', 'sortOrder'];

const journalQueryMap: Partial<Record<JournalFilterKey, string>> = journalFields.reduce((queryMap, field) => {
	const entry = journalFilterConfig[field];
	if (!entry.queryKey) return queryMap;
	return { ...queryMap, [field]: entry.queryKey };
}, {} as Partial<Record<JournalFilterKey, string>>);

const journalReverseMap: Record<string, JournalFilterKey> = journalFields.reduce((reverseMap, field) => {
	const queryKey = journalQueryMap[field] ?? field;
	const aliases = journalFilterConfig[field].aliases ?? [];

	return {
		...reverseMap,
		[queryKey]: field,
		...aliases.reduce<Record<string, JournalFilterKey>>((aliasMap, alias) => ({ ...aliasMap, [alias]: field }), {}),
	};
}, {} as Record<string, JournalFilterKey>);

const journalFilterUrlMapping = {
	fields: journalFields,
	queryMap: journalQueryMap,
	reverseMap: journalReverseMap,
} as const;

type TypeToggle = {
	label: string;
	className: string;
	nextType: string;
};

const typeToggleMap: Record<string, TypeToggle> = {
	'': { label: 'Taken', className: 'border-rose-300/70 bg-rose-100/60 text-rose-800 hover:bg-rose-200/70', nextType: 'TAKEN' },
	TAKEN: { label: 'Rejected', className: 'border-violet-300/70 bg-violet-100/60 text-violet-800 hover:bg-violet-200/70', nextType: 'REJECTED' },
	REJECTED: { label: 'All', className: 'border-slate-300/70 bg-slate-100/70 text-slate-700 hover:bg-slate-200/80', nextType: '' },
};

function createDefaultJournalFilters(): JournalFilterState {
	const defaults = journalFields.reduce<Record<JournalFilterKey, string>>((acc: Record<JournalFilterKey, string>, field: JournalFilterKey) => ({
		...acc,
		[field]: field === 'sortBy' ? 'created_at' : field === 'sortOrder' ? 'desc' : '',
	}), {} as Record<JournalFilterKey, string>);

	return {
		...defaults,
		clear() {},
		toQueryParams() {
			return { ...this } as Record<JournalFilterKey, string>;
		},
		hasActiveState() {
			return false;
		},
	} as JournalFilterState;
}

export function createJournalFilter(page: JournalPageData): JournalFilterState {
	const state = { ...createDefaultJournalFilters() } as JournalFilterState;

	state.clear = function clear(this: JournalFilterState) {
		Object.assign(this, createDefaultJournalFilters());
	};

	state.toQueryParams = function toQueryParams(this: JournalFilterState) {
		return { ...this };
	};

	state.hasActiveState = function hasActiveState(this: JournalFilterState) {
		return journalFields.some((field: JournalFilterKey) => {
			if (field === 'sortBy') return this.sortBy !== 'created_at';
			if (field === 'sortOrder') return this.sortOrder !== 'desc';
			return this[field] !== '';
		});
	};

	state.urlToFilter = function urlToFilter(this: JournalFilterState) {
		syncJournalUrlToFilter(this);
	};

	state.filterToUrl = function filterToUrl(this: JournalFilterState) {
		syncJournalFilterToUrl(this);
	};

	state.toggleType = function toggleType(this: JournalFilterState) {
		this.type = resolveTypeToggle(this.type).nextType;
		this.applyManualFilters();
	};

	state.typeToggleLabel = function typeToggleLabel(this: JournalFilterState) {
		return resolveTypeToggle(this.type).label;
	};

	state.typeToggleClass = function typeToggleClass(this: JournalFilterState) {
		return resolveTypeToggle(this.type).className;
	};

	state.onCreatedDateChange = function onCreatedDateChange(this: JournalFilterState) {
		this.createdBefore = this.createdAfter;
		this.applyManualFilters();
	};

	state.toggleSort = function toggleSort(this: JournalFilterState, field: string) {
		this.sortOrder = this.sortBy !== field ? 'asc' : this.sortOrder === 'asc' ? 'desc' : 'asc';
		this.sortBy = field;
		this.applyManualFilters();
	};

	state.applyManualFilters = function applyManualFilters() {
		page.table.applyManualFilters();
	};

	state.clearFilters = function clearFilters(this: JournalFilterState) {
		this.clear();
		this.applyManualFilters();
	};

	return state;
}

export function syncJournalFilterToUrl(filter: JournalFilterState) {
	syncStateToUrl(filter as unknown as JournalFilters, journalFilterUrlMapping);
}

export function syncJournalUrlToFilter(filter: JournalFilterState) {
	syncUrlToState(filter as unknown as JournalFilters, journalFilterUrlMapping);
}

export function resolveTypeToggle(currentType: string): TypeToggle {
	return typeToggleMap[currentType] ?? typeToggleMap[''];
}

import { syncStateToUrl, syncUrlToState } from '../shared/url_state';
import type { JournalFilterState } from './filter';
import { journalFilterUrlMapping } from './filter';

type TypeToggle = {
	label: string;
	className: string;
	nextType: string;
};

const typeToggleMap: Record<string, TypeToggle> = {
	'': {
		label: 'Taken',
		className: 'border-rose-300/70 bg-rose-100/60 text-rose-800 hover:bg-rose-200/70',
		nextType: 'TAKEN',
	},
	TAKEN: {
		label: 'Rejected',
		className: 'border-violet-300/70 bg-violet-100/60 text-violet-800 hover:bg-violet-200/70',
		nextType: 'REJECTED',
	},
	REJECTED: {
		label: 'All',
		className: 'border-slate-300/70 bg-slate-100/70 text-slate-700 hover:bg-slate-200/80',
		nextType: '',
	},
};

export function resolveTypeToggle(currentType: string): TypeToggle {
	return typeToggleMap[currentType] ?? typeToggleMap[''];
}

export function nextSortOrder(currentSortBy: string, currentSortOrder: string, nextField: string): string {
	if (currentSortBy !== nextField) {
		return 'asc';
	}

	return currentSortOrder === 'asc' ? 'desc' : 'asc';
}

type FilterActionDeps = {
	filter: JournalFilterState;
	applyManualFilters: () => void;
};

export function createFilterActions(deps: FilterActionDeps) {
	const { filter, applyManualFilters } = deps;

	return {
		urlToFilter() {
			syncUrlToState(filter, journalFilterUrlMapping);
		},
		filterToUrl() {
			syncStateToUrl(filter, journalFilterUrlMapping);
		},
		toggleType() {
			filter.type = resolveTypeToggle(filter.type).nextType;
			applyManualFilters();
		},
		onCreatedDateChange() {
			filter.createdBefore = filter.createdAfter;
			applyManualFilters();
		},
		toggleSort(field: string) {
			filter.sortOrder = nextSortOrder(filter.sortBy, filter.sortOrder, field);
			filter.sortBy = field;
			applyManualFilters();
		},
		clearFilters() {
			filter.clear();
			applyManualFilters();
		},
	};
}

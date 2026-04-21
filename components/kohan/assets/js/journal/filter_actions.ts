import { syncStateToUrl, syncUrlToState } from '../shared/url_state';
import type { JournalFilterState } from './filter';
import { journalFilterUrlMapping } from './filter';

type TypeToggle = {
	label: string;
	className: string;
	nextType: string;
};

const takenToggle: TypeToggle = {
	label: 'Rejected',
	className: '!border-emerald-300 !bg-emerald-200 !text-emerald-800',
	nextType: 'REJECTED',
};

const rejectedToggle: TypeToggle = {
	label: 'Taken',
	className: '!border-rose-300 !bg-rose-200 !text-rose-800',
	nextType: 'TAKEN',
};

export function resolveTypeToggle(type: string): TypeToggle {
	return type === 'TAKEN' ? takenToggle : rejectedToggle;
}

function nextSortOrder(sortBy: string, sortOrder: string, field: string): string {
	if (sortBy !== field) {
		return 'asc';
	}

	return sortOrder === 'asc' ? 'desc' : 'asc';
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

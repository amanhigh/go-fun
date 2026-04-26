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

export function createFilterActions() {
	return {
		urlToFilter(this: any) {
			syncUrlToState(this.filter as JournalFilterState, journalFilterUrlMapping);
		},
		filterToUrl(this: any) {
			syncStateToUrl(this.filter as JournalFilterState, journalFilterUrlMapping);
		},
		toggleType(this: any) {
			this.filter.type = resolveTypeToggle(this.filter.type).nextType;
			this.applyManualFilters();
		},
		onCreatedDateChange(this: any) {
			this.filter.createdBefore = this.filter.createdAfter;
			this.applyManualFilters();
		},
		toggleSort(this: any, field: string) {
			this.filter.sortOrder = nextSortOrder(this.filter.sortBy, this.filter.sortOrder, field);
			this.filter.sortBy = field;
			this.applyManualFilters();
		},
		clearFilters(this: any) {
			this.filter.clear();
			this.applyManualFilters();
		},
	};
}

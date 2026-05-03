import { syncStateToUrl, syncUrlToState } from '../../../shared/url_state';
import { journalFields, journalFilterUrlMapping, type JournalFilterKey, type JournalFilters, type JournalFilterState } from '../../../types/journal_list_state';

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

export function createJournalFilter(): JournalFilterState {
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

export function createFilterActions() {
	return {
		urlToFilter(this: any) { syncJournalUrlToFilter(this.filter); },
		filterToUrl(this: any) { syncJournalFilterToUrl(this.filter); },
		toggleType(this: any) { this.filter.type = resolveTypeToggle(this.filter.type).nextType; this.applyManualFilters(); },
		onCreatedDateChange(this: any) { this.filter.createdBefore = this.filter.createdAfter; this.applyManualFilters(); },
		toggleSort(this: any, field: string) { this.filter.sortOrder = this.filter.sortBy !== field ? 'asc' : this.filter.sortOrder === 'asc' ? 'desc' : 'asc'; this.filter.sortBy = field; this.applyManualFilters(); },
		clearFilters(this: any) { this.filter.clear(); this.applyManualFilters(); },
	};
}

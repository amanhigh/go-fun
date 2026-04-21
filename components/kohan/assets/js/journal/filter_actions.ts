import type { JournalFilterState } from './filter_state';

type TypeToggleState = {
	label: string;
	buttonClass: string;
	nextType: string;
};

const takenToggleState: TypeToggleState = {
	label: 'Rejected',
	buttonClass: '!border-emerald-300 !bg-emerald-200 !text-emerald-800',
	nextType: 'REJECTED',
};

const rejectedToggleState: TypeToggleState = {
	label: 'Taken',
	buttonClass: '!border-rose-300 !bg-rose-200 !text-rose-800',
	nextType: 'TAKEN',
};

export function getJournalTypeToggleState(type: string): TypeToggleState {
	return type === 'TAKEN' ? takenToggleState : rejectedToggleState;
}

export function createFilterActions(filter: JournalFilterState) {
	return {
		toggleTypeFilter(this: any) {
			this.filter.type = getJournalTypeToggleState(this.filter.type).nextType;
			this.applyManualFilters();
		},
		applyFilters(this: any) {
			this.pagination.resetPage();
			this.filterToUrl();
			void this.loadJournals();
		},
		onCreatedDateChange(this: any) {
			this.filter.createdBefore = this.filter.createdAfter;
			this.applyManualFilters();
		},
		applyManualFilters(this: any) {
			this.clearActiveReviewPreset();
			this.applyFilters();
		},
		toggleSort(this: any, field: string) {
			this.filter.sortOrder = this.filter.sortBy !== field
				? 'asc'
				: this.filter.sortOrder === 'asc' ? 'desc' : 'asc';
			this.filter.sortBy = field;
			this.applyManualFilters();
		},
		clearFilters(this: any) {
			this.filter.clear();
			this.applyManualFilters();
		},
	};
}

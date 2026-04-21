import type { JournalFilterState } from './filter_state';

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

export function createFilterActions(filter: JournalFilterState) {
	return {
		toggleType(this: any) {
			this.filter.type = resolveTypeToggle(this.filter.type).nextType;
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

import type { JournalFilterState } from './filter_state';

export function createFilterActions(filter: JournalFilterState) {
	return {
		toggleTypeFilter(this: any) {
			this.filter.type = this.filter.type === 'TAKEN' ? 'REJECTED' : 'TAKEN';
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

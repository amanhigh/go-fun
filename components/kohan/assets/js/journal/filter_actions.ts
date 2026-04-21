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
			this.clearActiveReviewPreset();
			this.filter.createdBefore = this.filter.createdAfter;
			this.applyFilters();
		},
		applyManualFilters(this: any) {
			this.clearActiveReviewPreset();
			this.applyFilters();
		},
		toggleSort(this: any, field: string) {
			this.clearActiveReviewPreset();
			this.filter.sortOrder = this.filter.sortBy !== field
				? 'asc'
				: this.filter.sortOrder === 'asc' ? 'desc' : 'asc';
			this.filter.sortBy = field;
			this.applyFilters();
		},
		clearFilters(this: any) {
			this.clearActiveReviewPreset();
			this.filter.clear();
			this.applyFilters();
		},
	};
}

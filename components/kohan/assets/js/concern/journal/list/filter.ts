import { JournalSortBy, JournalSortOrder } from '../../../types/api/journal/enums';
import { ReviewedFilter } from '../../../types/api/journal/request';
import type { JournalFilterValues, DatePresetName, JournalFilterConcern, JournalPageProvider } from '../../../types/journal/list';

const journalFilterDefaults: JournalFilterValues = {
	ticker: '',
	type: '',
	status: '',
	sequence: '',
	createdAfter: '',
	createdBefore: '',
	reviewed: ReviewedFilter.ALL,
	sortBy: JournalSortBy.CREATED_AT,
	sortOrder: JournalSortOrder.DESC,
};

export function NewFilterConcern(pg: JournalPageProvider): JournalFilterConcern {
	return {
		...journalFilterDefaults,
		datePreset: '' as DatePresetName,
		clear() {
			Object.assign(this, journalFilterDefaults);
			this.datePreset = '';
		},
		hasActiveState() {
			if (this.datePreset !== '') return true;
			return (Object.keys(journalFilterDefaults) as (keyof JournalFilterValues)[]).some((field) => this[field] !== journalFilterDefaults[field]);
		},
		applyFilters() {
			pg().pagination.resetPage();
			pg().filterUrl.filterToUrl();
			void pg().table.load();
		},
		applyCreatedDate(createdAt: string) {
			const date = pg().present.date.humanDate(new Date(createdAt));
			this.createdAfter = date;
			this.createdBefore = date;
			this.datePreset = '';
			pg().presets.clearActiveReviewPreset();
			this.applyFilters();
		},
		applyManualFilters() {
			pg().presets.clearActiveReviewPreset();
			this.datePreset = '';
			this.applyFilters();
		},
		toggleSort(field: JournalSortBy) {
			if (this.sortBy !== field) {
				this.sortOrder = JournalSortOrder.ASC;
			} else {
				this.sortOrder = this.sortOrder === JournalSortOrder.ASC ? JournalSortOrder.DESC : JournalSortOrder.ASC;
			}
			this.sortBy = field;
			this.applyManualFilters();
		},
		clearFilters() {
			this.clear();
			this.applyManualFilters();
		},
	} as JournalFilterConcern;
}

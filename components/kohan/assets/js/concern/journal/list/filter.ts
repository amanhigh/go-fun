import { JournalType, JournalSortBy, JournalSortOrder } from '../../../types/api/journal/enums';
import { ReviewedFilter } from '../../../types/api/journal/request';
import type { JournalFilterValues, DatePresetName, JournalFilterConcern, JournalPageProvider } from '../../../types/journal/list';

type TypeToggle = {
	label: string;
	className: string;
	nextType: JournalType | '';
};

const typeToggleMap: Record<string, TypeToggle> = {
	'': { label: 'Taken', className: 'journal-type-toggle-taken', nextType: JournalType.TAKEN },
	[JournalType.TAKEN]: { label: 'Rejected', className: 'journal-type-toggle-rejected', nextType: JournalType.REJECTED },
	[JournalType.REJECTED]: { label: 'All', className: 'journal-type-toggle-all', nextType: '' },
};

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
		typeToggle() {
			return resolveTypeToggle(this.type);
		},
		applyFilters() {
			pg().pagination.resetPage();
			pg().filterUrl.filterToUrl();
			void pg().table.loadJournals();
		},
		applyManualFilters() {
			pg().presets.clearActiveReviewPreset();
			this.datePreset = '';
			this.applyFilters();
		},
		toggleType() {
			this.type = this.typeToggle().nextType;
			this.applyManualFilters();
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

export function resolveTypeToggle(currentType: JournalType | ''): TypeToggle {
	if (currentType === JournalType.TAKEN || currentType === JournalType.REJECTED) {
		return typeToggleMap[currentType];
	}
	return typeToggleMap[''];
}

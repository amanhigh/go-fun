import type { JournalFilterKey } from '../../../types/journal_api';
import type { DatePresetName, JournalFilterConcern, JournalPageProvider } from '../../../types/journal_list_concern';

type TypeToggle = {
	label: string;
	className: string;
	nextType: TypeFilterValue;
};

type SortField = 'ticker' | 'sequence' | 'created_at';

type TypeFilterValue = '' | 'TAKEN' | 'REJECTED';

const typeToggleMap: Record<TypeFilterValue, TypeToggle> = {
	'': { label: 'Taken', className: 'journal-type-toggle-taken', nextType: 'TAKEN' },
	TAKEN: { label: 'Rejected', className: 'journal-type-toggle-rejected', nextType: 'REJECTED' },
	REJECTED: { label: 'All', className: 'journal-type-toggle-all', nextType: '' },
};

const journalFilterDefaults: Record<JournalFilterKey, string> = {
	ticker: '',
	type: '',
	status: '',
	sequence: '',
	createdAfter: '',
	createdBefore: '',
	reviewed: '',
	sortBy: 'created_at',
	sortOrder: 'desc',
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
			return Object.entries(journalFilterDefaults).some(([field, defaultValue]) => this[field as JournalFilterKey] !== defaultValue);
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
		toggleSort(field: SortField) {
			if (this.sortBy !== field) {
				this.sortOrder = 'asc';
			} else {
				this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
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

export function resolveTypeToggle(currentType: string): TypeToggle {
	if (currentType === 'TAKEN' || currentType === 'REJECTED') {
		return typeToggleMap[currentType];
	}
	return typeToggleMap[''];
}

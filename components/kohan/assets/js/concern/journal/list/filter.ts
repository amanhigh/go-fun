import type { JournalFilterKey } from '../../../types/journal_api';
import type { DatePresetName, JournalFilterConcern, JournalPageProvider } from '../../../types/journal_list_concern';

type TypeToggle = {
	label: string;
	className: string;
	nextType: string;
};

type SortField = 'ticker' | 'sequence' | 'created_at';

const typeToggleMap: Record<string, TypeToggle> = {
	'': { label: 'Taken', className: 'border-rose-300/70 bg-rose-100/60 text-rose-800 hover:bg-rose-200/70', nextType: 'TAKEN' },
	TAKEN: { label: 'Rejected', className: 'border-violet-300/70 bg-violet-100/60 text-violet-800 hover:bg-violet-200/70', nextType: 'REJECTED' },
	REJECTED: { label: 'All', className: 'border-slate-300/70 bg-slate-100/70 text-slate-700 hover:bg-slate-200/80', nextType: '' },
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
			return Object.entries(journalFilterDefaults).some(([field, defaultValue]) => this[field as JournalFilterKey] !== defaultValue);
		},
		typeToggle() {
			return resolveTypeToggle(this.type);
		},
		toggleType() {
			this.type = this.typeToggle().nextType;
			pg().table.applyManualFilters();
		},
		toggleSort(field: SortField) {
			this.sortOrder = this.sortBy !== field ? 'asc' : this.sortOrder === 'asc' ? 'desc' : 'asc';
			this.sortBy = field;
			pg().table.applyManualFilters();
		},
		applyManualFilters() {
			pg().table.applyManualFilters();
		},
		clearFilters() {
			this.clear();
			pg().table.applyManualFilters();
		},
	} as JournalFilterConcern;
}

export function resolveTypeToggle(currentType: string): TypeToggle {
	return typeToggleMap[currentType] ?? typeToggleMap[''];
}

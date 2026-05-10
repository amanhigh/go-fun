import { JournalType, JournalStatus, JournalSortBy, JournalSortOrder } from '../../../types/api/journal/enums';
import { ReviewedFilter } from '../../../types/api/journal/request';
import type { JournalFilterValues, DatePresetName, JournalFilterConcern, JournalPageProvider } from '../../../types/journal/list';

type TypeToggle = {
	label: string;
	className: string;
	nextType: JournalType | '';
};

type StatusToggle = {
	label: string;
	className: string;
	nextStatus: JournalStatus | '';
};

const allToggleSpec: { label: string; className: string } = {
	label: 'All',
	className: 'journal-display-default',
};

const typeTransitionMap: Record<string, { nextType: JournalType | '' }> = {
	'': { nextType: JournalType.TAKEN },
	[JournalType.TAKEN]: { nextType: JournalType.REJECTED },
	[JournalType.REJECTED]: { nextType: '' },
};

const statusTransitionMap: Record<string, { nextStatus: JournalStatus | '' }> = {
	'': { nextStatus: JournalStatus.SET },
	[JournalStatus.SET]: { nextStatus: JournalStatus.RUNNING },
	[JournalStatus.RUNNING]: { nextStatus: '' },
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
			return resolveTypeToggle(pg, this.type);
		},
		statusToggle() {
			return resolveStatusToggle(pg, this.status);
		},
		toggleStatus() {
			this.status = this.statusToggle().nextStatus;
			this.applyManualFilters();
		},
		applyFilters() {
			pg().pagination.resetPage();
			pg().filterUrl.filterToUrl();
			void pg().table.load();
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

export function resolveTypeToggle(pg: JournalPageProvider, currentType: JournalType | ''): TypeToggle {
	const transition = typeTransitionMap[currentType] ?? typeTransitionMap[''];
	if (transition.nextType === '') {
		return { label: allToggleSpec.label, className: allToggleSpec.className, nextType: '' };
	}
	return {
		label: pg().present.type.label(transition.nextType),
		className: pg().present.type.spec(transition.nextType).class,
		nextType: transition.nextType,
	};
}

export function resolveStatusToggle(pg: JournalPageProvider, currentStatus: JournalStatus | ''): StatusToggle {
	const transition = statusTransitionMap[currentStatus] ?? statusTransitionMap[''];
	if (transition.nextStatus === '') {
		return { label: allToggleSpec.label, className: allToggleSpec.className, nextStatus: '' };
	}
	return {
		label: pg().present.status.label(transition.nextStatus),
		className: pg().present.status.spec(transition.nextStatus).class,
		nextStatus: transition.nextStatus,
	};
}

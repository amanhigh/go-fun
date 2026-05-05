import type { DatePresetName, JournalPageData, JournalPageProvider, PresetConcern, ReviewPreset } from '../../../types/journal_list_concern';
import { formatDateInputValue } from '../../../shared/date';

const monthLabels = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
const reviewPresetOffsets = [-2, -1, 0, 1, 2] as const;
const reviewPresetAnchorMonthOffset = 9;

type ReviewPresetFilter = Pick<JournalPageData['filter'], 'createdAfter' | 'createdBefore' | 'reviewed'>;

const datePresetMap: Record<Exclude<DatePresetName, ''>, number> = {
	today: 0,
	last7: 7,
	last30: 30,
};

export function buildDatePresetRange(preset: DatePresetName, today: Date = new Date()): { createdAfter: string; createdBefore: string } {
	if (!preset) return { createdAfter: '', createdBefore: '' };
	const days = datePresetMap[preset];
	const endDate = formatDateInputValue(today);
	const startDate = new Date(today);
	startDate.setDate(today.getDate() - days);

	return {
		createdAfter: formatDateInputValue(startDate),
		createdBefore: endDate,
	};
}

export function syncDatePreset(filter: JournalPageData['filter']) {
	if (filter.datePreset) {
		const range = buildDatePresetRange(filter.datePreset);
		filter.createdAfter = range.createdAfter;
		filter.createdBefore = range.createdBefore;
	}
}

function buildReviewPresetAnchorDate(today: Date = new Date()): Date {
	return new Date(today.getFullYear(), today.getMonth() - reviewPresetAnchorMonthOffset, 1);
}

function resolveReviewPresetClass(reviewPreset: ReviewPreset, activeReviewPreset: string): string {
	if (activeReviewPreset === reviewPreset.label) return reviewPresetActiveClass;
	if (reviewPreset.isAnchor) return reviewPresetAnchorClass;
	return reviewPresetBaseClass;
}

export function buildReviewPresetList(): ReviewPreset[] {
	const today = new Date();
	const anchorDate = buildReviewPresetAnchorDate(today);
	return reviewPresetOffsets.map((offset) => {
		const monthDate = new Date(anchorDate.getFullYear(), anchorDate.getMonth() + offset, 1);
		const createdAfterDate = new Date(monthDate.getFullYear(), monthDate.getMonth(), 1);
		const createdBeforeDate = new Date(monthDate.getFullYear(), monthDate.getMonth() + 1, 0);
		return {
			isAnchor: offset === 0,
			label: `${monthLabels[monthDate.getMonth()]}-${String(monthDate.getFullYear() % 100).padStart(2, '0')}`,
			createdAfter: formatDateInputValue(createdAfterDate),
			createdBefore: formatDateInputValue(createdBeforeDate),
		};
	});
}

export function findReviewPreset(reviewPresets: ReviewPreset[], filter: ReviewPresetFilter): ReviewPreset | undefined {
	return reviewPresets.find((reviewPreset) => (
		reviewPreset.createdAfter === filter.createdAfter
		&& reviewPreset.createdBefore === filter.createdBefore
		&& filter.reviewed === 'false'
	));
}

const reviewPresetBaseClass = 'border-cyan-200/70 bg-white/80 text-cyan-800 hover:bg-cyan-50';
const reviewPresetAnchorClass = 'border-2 border-amber-200 bg-white/80 text-cyan-800';
const reviewPresetActiveClass = 'border-amber-300 bg-amber-100/90 text-amber-950 hover:bg-amber-100';

function applyPresetChanges(pg: JournalPageProvider, presets: PresetConcern, activeReviewPreset: string, mutate: () => void) {
	pg().filter.clear();
	mutate();
	presets.activeReviewPreset = activeReviewPreset;
	pg().table.applyFilters();
}

export function newPresetConcern(pg: JournalPageProvider): PresetConcern {
	const presets: PresetConcern = {
		reviewPresets: buildReviewPresetList(),
		activeReviewPreset: '',
		clearActiveReviewPreset() { presets.activeReviewPreset = ''; },
		syncActiveReviewPreset() {
			presets.activeReviewPreset = findReviewPreset(presets.reviewPresets, pg().filter)?.label ?? '';
		},
		syncDatePreset() {
			syncDatePreset(pg().filter);
		},
		reviewPresetClass(reviewPreset: ReviewPreset) { return resolveReviewPresetClass(reviewPreset, presets.activeReviewPreset); },
		applyCreatedPreset(preset: DatePresetName) {
			applyPresetChanges(pg, presets, '', () => {
				pg().filter.datePreset = preset;
				const range = buildDatePresetRange(preset);
				pg().filter.createdAfter = range.createdAfter;
				pg().filter.createdBefore = range.createdBefore;
			});
		},
		applyReviewPreset(reviewPreset: ReviewPreset) {
			applyPresetChanges(pg, presets, reviewPreset.label, () => {
				pg().filter.createdAfter = reviewPreset.createdAfter;
				pg().filter.createdBefore = reviewPreset.createdBefore;
				pg().filter.reviewed = 'false';
			});
		},
	};

	return presets;
}

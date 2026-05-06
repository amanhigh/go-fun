import type { DatePresetName, JournalPageData, JournalPageProvider, NonEmptyDatePresetName, PresetConcern, ReviewPreset } from '../../../types/journal_list_concern';
import { formatDateInputValue } from '../../../shared/date';

const monthLabels = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
const reviewPresetMonthOffsets = [-11, -10, -9, -8, -7] as const;
const reviewPresetAnchorOffset = -9;

const pendingReviewValue = 'false';

type ReviewPresetFilter = Pick<JournalPageData['filter'], 'createdAfter' | 'createdBefore' | 'reviewed'>;

const datePresetMap: Record<NonEmptyDatePresetName, number> = {
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
	if (!filter.datePreset) return;
	const range = buildDatePresetRange(filter.datePreset);
	filter.createdAfter = range.createdAfter;
	filter.createdBefore = range.createdBefore;
}

export function buildReviewPresetList(): ReviewPreset[] {
	const today = new Date();
	return reviewPresetMonthOffsets.map((offset) => {
		const monthDate = new Date(today.getFullYear(), today.getMonth() + offset, 1);
		const createdAfterDate = new Date(monthDate.getFullYear(), monthDate.getMonth(), 1);
		const createdBeforeDate = new Date(monthDate.getFullYear(), monthDate.getMonth() + 1, 0);
		return {
			isAnchor: offset === reviewPresetAnchorOffset,
			label: `${monthLabels[monthDate.getMonth()]}-${String(monthDate.getFullYear() % 100).padStart(2, '0')}`,
			createdAfter: formatDateInputValue(createdAfterDate),
			createdBefore: formatDateInputValue(createdBeforeDate),
		};
	});
}

function resolveReviewPresetClass(reviewPreset: ReviewPreset, activeReviewPreset: string): string {
	if (activeReviewPreset === reviewPreset.label) return reviewPresetActiveClass;
	if (reviewPreset.isAnchor) return reviewPresetAnchorClass;
	return reviewPresetBaseClass;
}

export function findReviewPreset(reviewPresets: ReviewPreset[], filter: ReviewPresetFilter): ReviewPreset | undefined {
	return reviewPresets.find((reviewPreset) => (
		reviewPreset.createdAfter === filter.createdAfter
		&& reviewPreset.createdBefore === filter.createdBefore
		&& filter.reviewed === pendingReviewValue
	));
}

const reviewPresetBaseClass = 'journal-review-preset-base';
const reviewPresetAnchorClass = 'journal-review-preset-anchor';
const reviewPresetActiveClass = 'journal-review-preset-active';

function applyPresetChanges(pg: JournalPageProvider, presets: PresetConcern, activeReviewPreset: string, mutate: () => void) {
	pg().filter.clear();
	mutate();
	presets.activeReviewPreset = activeReviewPreset;
	pg().filter.applyFilters();
}

export function NewPresetConcern(pg: JournalPageProvider): PresetConcern {
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
		applyCreatedPreset(preset: NonEmptyDatePresetName) {
			const filter = pg().filter;
			filter.clear();
			filter.datePreset = preset;
			syncDatePreset(filter);
			presets.activeReviewPreset = '';
			filter.applyFilters();
		},
		applyReviewPreset(reviewPreset: ReviewPreset) {
			const filter = pg().filter;
			filter.clear();
			filter.createdAfter = reviewPreset.createdAfter;
			filter.createdBefore = reviewPreset.createdBefore;
			filter.reviewed = pendingReviewValue;
			presets.activeReviewPreset = reviewPreset.label;
			filter.applyFilters();
		},
	};

	return presets;
}

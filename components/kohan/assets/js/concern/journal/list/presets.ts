import type { JournalPageData, PresetState, ReviewPreset } from '../../../types/journal_list_state';
import { formatDateInputValue } from '../../../shared/date';

const monthLabels = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
const reviewPresetOffsets = [-2, -1, 0, 1, 2] as const;

export type DatePreset = 'today' | 'last7' | 'last30';

const datePresetMap: Record<DatePreset, number> = {
	today: 0,
	last7: 7,
	last30: 30,
};

export function buildDatePresetRange(preset: DatePreset, today: Date = new Date()): { createdAfter: string; createdBefore: string } {
	const days = datePresetMap[preset] ?? datePresetMap.last7;
	const endDate = formatDateInputValue(today);
	const startDate = new Date(today);
	startDate.setDate(today.getDate() - days);

	return {
		createdAfter: formatDateInputValue(startDate),
		createdBefore: endDate,
	};
}

export function buildReviewPresetList(): ReviewPreset[] {
	const today = new Date();
	const anchorDate = new Date(today.getFullYear(), today.getMonth() - 9, 1);
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

export function findReviewPreset(reviewPresets: ReviewPreset[], filter: JournalPageData['filter']): ReviewPreset | undefined {
	return reviewPresets.find((reviewPreset) => (
		reviewPreset.createdAfter === filter.createdAfter
		&& reviewPreset.createdBefore === filter.createdBefore
		&& filter.reviewed === 'false'
	));
}

const reviewPresetBaseClass = 'border-cyan-200/70 bg-white/80 text-cyan-800 hover:bg-cyan-50';
const reviewPresetAnchorClass = 'border-2 border-amber-200 bg-white/80 text-cyan-800';
const reviewPresetActiveClass = 'border-amber-300 bg-amber-100/90 text-amber-950 hover:bg-amber-100';

export function createPresetConcern(page: JournalPageData): PresetState {
	const presets: PresetState = {
		reviewPresets: buildReviewPresetList(),
		activeReviewPreset: '',
		clearActiveReviewPreset() { presets.activeReviewPreset = ''; },
		syncActiveReviewPreset() { presets.activeReviewPreset = findReviewPreset(presets.reviewPresets, page.filter)?.label ?? ''; },
		reviewPresetClass(reviewPreset: ReviewPreset) { return presets.activeReviewPreset === reviewPreset.label ? reviewPresetActiveClass : reviewPreset.isAnchor ? reviewPresetAnchorClass : reviewPresetBaseClass; },
		applyCreatedPreset(preset: DatePreset) {
			page.filter.clear();
			const range = buildDatePresetRange(preset);
			page.filter.createdAfter = range.createdAfter;
			page.filter.createdBefore = range.createdBefore;
			presets.activeReviewPreset = '';
			page.table.applyFilters();
		},
		applyReviewPreset(reviewPreset: ReviewPreset) {
			page.filter.clear();
			page.filter.createdAfter = reviewPreset.createdAfter;
			page.filter.createdBefore = reviewPreset.createdBefore;
			page.filter.reviewed = 'false';
			presets.activeReviewPreset = reviewPreset.label;
			page.table.applyFilters();
		},
	};

	return presets;
}

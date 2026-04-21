import type { JournalFilterState } from './filter';
import { formatDateInputValue } from '../shared/date';

const monthLabels = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
const reviewPresetOffsets = [-2, -1, 0, 1, 2] as const;

export type ReviewPreset = {
	isAnchor: boolean;
	label: string;
	createdAfter: string;
	createdBefore: string;
};

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

export function findReviewPreset(reviewPresets: ReviewPreset[], filter: JournalFilterState): ReviewPreset | undefined {
	return reviewPresets.find((reviewPreset) => (
		reviewPreset.createdAfter === filter.createdAfter
		&& reviewPreset.createdBefore === filter.createdBefore
		&& filter.reviewed === 'false'
	));
}

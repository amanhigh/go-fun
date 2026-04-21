import { formatDateInputValue } from '../shared/date';
import type { JournalFilterState } from './filter';

const monthLabels = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
const reviewPresetOffsets = [-2, -1, 0, 1, 2] as const;

const reviewPresetBaseClass = 'border-cyan-200/70 bg-white/80 text-cyan-800 hover:bg-cyan-50';
const reviewPresetAnchorClass = 'border-2 border-amber-200 bg-white/80 text-cyan-800';
const reviewPresetActiveClass = 'border-amber-300 bg-amber-100/90 text-amber-950 hover:bg-amber-100';

function formatReviewPresetLabel(date: Date): string {
	return `${monthLabels[date.getMonth()]}-${String(date.getFullYear() % 100).padStart(2, '0')}`;
}

export type ReviewPreset = {
	isAnchor: boolean;
	label: string;
	createdAfter: string;
	createdBefore: string;
};

export type CreatedPreset = 'today' | 'last7' | 'last30';

const createdPresetMap: Record<CreatedPreset, number> = {
	today: 0,
	last7: 7,
	last30: 30,
};

function makeCreatedRange(days: number, today: Date): { createdAfter: string; createdBefore: string } {
	const endDate = formatDateInputValue(today);
	if (days === 0) {
		return {
			createdAfter: endDate,
			createdBefore: endDate,
		};
	}

	const startDate = new Date(today);
	startDate.setDate(today.getDate() - days);
	return {
		createdAfter: formatDateInputValue(startDate),
		createdBefore: endDate,
	};
}

export function getCreatedPresetRange(preset: CreatedPreset, today: Date = new Date()): { createdAfter: string; createdBefore: string } {
	const days = createdPresetMap[preset] ?? createdPresetMap.last7;
	return makeCreatedRange(days, today);
}

function createReviewPreset(monthDate: Date, isAnchor: boolean): ReviewPreset {
	const createdAfter = new Date(monthDate.getFullYear(), monthDate.getMonth(), 1);
	const createdBefore = new Date(monthDate.getFullYear(), monthDate.getMonth() + 1, 0);

	return {
		isAnchor,
		label: formatReviewPresetLabel(monthDate),
		createdAfter: formatDateInputValue(createdAfter),
		createdBefore: formatDateInputValue(createdBefore),
	};
}

export function createReviewPresets(): ReviewPreset[] {
	const today = new Date();
	const anchorDate = new Date(today.getFullYear(), today.getMonth() - 9, 1);
	return reviewPresetOffsets.map((offset) => {
		const monthDate = new Date(anchorDate.getFullYear(), anchorDate.getMonth() + offset, 1);
		return createReviewPreset(monthDate, offset === 0);
	});
}

export function getReviewPresetClass(activeReviewPreset: string, reviewPreset: ReviewPreset): string {
	if (activeReviewPreset === reviewPreset.label) {
		return reviewPresetActiveClass;
	}
	return reviewPreset.isAnchor ? reviewPresetAnchorClass : reviewPresetBaseClass;
}

export function findMatchingReviewPreset(reviewPresets: ReviewPreset[], filter: JournalFilterState): ReviewPreset | undefined {
	return reviewPresets.find((reviewPreset) => (
		reviewPreset.createdAfter === filter.createdAfter
		&& reviewPreset.createdBefore === filter.createdBefore
		&& filter.reviewed === 'false'
	));
}

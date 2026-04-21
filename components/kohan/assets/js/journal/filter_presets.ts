import { formatDateInputValue } from '../shared/date';

const reviewPresetBaseClass = 'border-cyan-200/70 bg-white/80 text-cyan-800 hover:bg-cyan-50';
const reviewPresetAnchorClass = 'border-2 border-amber-200 bg-white/80 text-cyan-800';
const reviewPresetActiveClass = 'border-amber-300 bg-amber-100/90 text-amber-950 hover:bg-amber-100';

function formatReviewPresetLabel(date: Date): string {
	const monthLabels = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
	return `${monthLabels[date.getMonth()]}-${String(date.getFullYear() % 100).padStart(2, '0')}`;
}

export type ReviewPreset = {
	isAnchor: boolean;
	label: string;
	createdAfter: string;
	createdBefore: string;
};

export function createReviewPresets(): ReviewPreset[] {
	const today = new Date();
	const anchorDate = new Date(today.getFullYear(), today.getMonth() - 9, 1);
	return [-2, -1, 0, 1, 2].map((offset) => {
		const monthDate = new Date(anchorDate.getFullYear(), anchorDate.getMonth() + offset, 1);
		const createdAfter = new Date(monthDate.getFullYear(), monthDate.getMonth(), 1);
		const createdBefore = new Date(monthDate.getFullYear(), monthDate.getMonth() + 1, 0);
		return {
			isAnchor: offset === 0,
			label: formatReviewPresetLabel(monthDate),
			createdAfter: formatDateInputValue(createdAfter),
			createdBefore: formatDateInputValue(createdBefore),
		};
	});
}

export function reviewPresetButtonClass(activeReviewPreset: string, reviewPreset: ReviewPreset): string {
	if (activeReviewPreset === reviewPreset.label) {
		return reviewPresetActiveClass;
	}
	return reviewPreset.isAnchor ? reviewPresetAnchorClass : reviewPresetBaseClass;
}

export function createFilterPresetActions() {
	return {
		currentReviewPresetLabel(this: any) {
			return this.activeReviewPreset;
		},
		syncActiveReviewPreset(this: any) {
			const matchingPreset = this.reviewPresets.find((reviewPreset: ReviewPreset) => (
				reviewPreset.createdAfter === this.filter.createdAfter
				&& reviewPreset.createdBefore === this.filter.createdBefore
				&& this.filter.reviewed === 'false'
			));
			this.activeReviewPreset = matchingPreset?.label ?? '';
		},
		clearActiveReviewPreset(this: any) {
			this.activeReviewPreset = '';
		},
		reviewPresetButtonClass(this: any, reviewPreset: ReviewPreset) {
			return reviewPresetButtonClass(this.activeReviewPreset, reviewPreset);
		},
		applyCreatedPreset(this: any, preset: string) {
			this.filter.clear();
			this.clearActiveReviewPreset();
			const today = new Date();
			const endDate = formatDateInputValue(today);
			const daysMap: Record<string, number> = { today: 0, last7: 7, last30: 30 };
			const days = daysMap[preset] ?? 7;
			if (days === 0) {
				this.filter.createdAfter = endDate;
				this.filter.createdBefore = endDate;
				this.applyFilters();
				return;
			}
			const startDate = new Date(today);
			startDate.setDate(today.getDate() - days);
			this.filter.createdAfter = formatDateInputValue(startDate);
			this.filter.createdBefore = endDate;
			this.applyFilters();
		},
		applyReviewPreset(this: any, reviewPreset: ReviewPreset) {
			this.filter.clear();
			this.filter.createdAfter = reviewPreset.createdAfter;
			this.filter.createdBefore = reviewPreset.createdBefore;
			this.filter.reviewed = 'false';
			this.activeReviewPreset = reviewPreset.label;
			this.applyFilters();
		},
	};
}

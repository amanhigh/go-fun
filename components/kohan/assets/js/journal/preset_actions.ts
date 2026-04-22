import { buildDatePresetRange, findReviewPreset, type DatePreset, type ReviewPreset } from './presets';

const reviewPresetBaseClass = 'border-cyan-200/70 bg-white/80 text-cyan-800 hover:bg-cyan-50';
const reviewPresetAnchorClass = 'border-2 border-amber-200 bg-white/80 text-cyan-800';
const reviewPresetActiveClass = 'border-amber-300 bg-amber-100/90 text-amber-950 hover:bg-amber-100';

export function createPresetActions() {
	return {
		clearActiveReviewPreset(this: any) {
			this.activeReviewPreset = '';
		},
		syncActiveReviewPreset(this: any) {
			const matchingPreset = findReviewPreset(this.reviewPresets, this.filter);
			this.activeReviewPreset = matchingPreset?.label ?? '';
		},
		reviewPresetClass(this: any, reviewPreset: ReviewPreset) {
			if (this.activeReviewPreset === reviewPreset.label) {
				return reviewPresetActiveClass;
			}

			return reviewPreset.isAnchor ? reviewPresetAnchorClass : reviewPresetBaseClass;
		},
		applyCreatedPreset(this: any, preset: DatePreset) {
			this.clearFilters();
			const range = buildDatePresetRange(preset);
			this.filter.createdAfter = range.createdAfter;
			this.filter.createdBefore = range.createdBefore;
			this.activeReviewPreset = '';
			this.applyFilters();
		},
		applyReviewPreset(this: any, reviewPreset: ReviewPreset) {
			this.clearFilters();
			this.filter.createdAfter = reviewPreset.createdAfter;
			this.filter.createdBefore = reviewPreset.createdBefore;
			this.filter.reviewed = 'false';
			this.activeReviewPreset = reviewPreset.label;
			this.applyFilters();
		},
	};
}

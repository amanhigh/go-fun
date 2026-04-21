import type { JournalFilterState } from './filter';
import type { JournalPageState } from './page_state';
import { buildDatePresetRange, findReviewPreset, type DatePreset, type ReviewPreset } from './presets';

const reviewPresetBaseClass = 'border-cyan-200/70 bg-white/80 text-cyan-800 hover:bg-cyan-50';
const reviewPresetAnchorClass = 'border-2 border-amber-200 bg-white/80 text-cyan-800';
const reviewPresetActiveClass = 'border-amber-300 bg-amber-100/90 text-amber-950 hover:bg-amber-100';

type PresetActionDeps = {
	filter: JournalFilterState;
	state: JournalPageState;
	applyFilters: () => void;
	clearFilters: () => void;
};

export function createPresetActions(deps: PresetActionDeps) {
	const { filter, state, applyFilters, clearFilters } = deps;

	return {
		clearActiveReviewPreset() {
			state.activeReviewPreset = '';
		},
		syncActiveReviewPreset() {
			const matchingPreset = findReviewPreset(state.reviewPresets, filter);
			state.activeReviewPreset = matchingPreset?.label ?? '';
		},
		reviewPresetClass(reviewPreset: ReviewPreset) {
			if (state.activeReviewPreset === reviewPreset.label) {
				return reviewPresetActiveClass;
			}

			return reviewPreset.isAnchor ? reviewPresetAnchorClass : reviewPresetBaseClass;
		},
		applyCreatedPreset(preset: DatePreset) {
			clearFilters();
			const range = buildDatePresetRange(preset);
			filter.createdAfter = range.createdAfter;
			filter.createdBefore = range.createdBefore;
			state.activeReviewPreset = '';
			applyFilters();
		},
		applyReviewPreset(reviewPreset: ReviewPreset) {
			clearFilters();
			filter.createdAfter = reviewPreset.createdAfter;
			filter.createdBefore = reviewPreset.createdBefore;
			filter.reviewed = 'false';
			state.activeReviewPreset = reviewPreset.label;
			applyFilters();
		},
	};
}

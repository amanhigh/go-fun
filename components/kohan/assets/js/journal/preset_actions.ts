import type { JournalFilterState } from './filter';
import type { JournalPageState } from './page_state';
import { findMatchingReviewPreset, getCreatedPresetRange, getReviewPresetClass, type CreatedPreset, type ReviewPreset } from './presets';

type PresetActionDeps = {
	filter: JournalFilterState;
	state: JournalPageState;
	applyFilters: () => void;
	clearFilters: () => void;
};

function setCreatedRange(filter: JournalFilterState, createdAfter: string, createdBefore: string) {
	filter.createdAfter = createdAfter;
	filter.createdBefore = createdBefore;
}

export function createPresetActions(deps: PresetActionDeps) {
	const { filter, state, applyFilters, clearFilters } = deps;

	function setActiveReviewPreset(label: string) {
		state.activeReviewPreset = label;
	}

	function clearActiveReviewPreset() {
		setActiveReviewPreset('');
	}

	function syncActiveReviewPreset() {
		const matchingPreset = findMatchingReviewPreset(state.reviewPresets, filter);
		setActiveReviewPreset(matchingPreset?.label ?? '');
	}

	return {
		clearActiveReviewPreset,
		syncActiveReviewPreset,
		reviewPresetClass(reviewPreset: ReviewPreset) {
			return getReviewPresetClass(state.activeReviewPreset, reviewPreset);
		},
		applyCreatedPreset(preset: CreatedPreset) {
			clearFilters();
			const range = getCreatedPresetRange(preset);
			setCreatedRange(filter, range.createdAfter, range.createdBefore);
			clearActiveReviewPreset();
			applyFilters();
		},
		applyReviewPreset(reviewPreset: ReviewPreset) {
			clearFilters();
			setCreatedRange(filter, reviewPreset.createdAfter, reviewPreset.createdBefore);
			filter.reviewed = 'false';
			setActiveReviewPreset(reviewPreset.label);
			applyFilters();
		},
	};
}

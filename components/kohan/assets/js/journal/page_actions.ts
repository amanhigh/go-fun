import type { Journal, JournalClient } from '../client/journal';
import type { JournalFilterState } from './filter_state';
import type { PaginationState } from './pagination';
import { findMatchingReviewPreset, getCreatedPresetRange, getReviewPresetClass, type CreatedPreset, type ReviewPreset } from './filter_presets';

type JournalPageState = {
	journals: Journal[];
	activeReviewPreset: string;
	requestCounter: number;
	loading: boolean;
	errorMessage: string;
};

type PageActionDeps = {
	client: JournalClient;
	filter: JournalFilterState;
	pagination: PaginationState;
	reviewPresets: ReviewPreset[];
	state: JournalPageState;
	filterToUrl: () => void;
	urlToFilter: () => void;
};

function applyCreatedRange(filter: JournalFilterState, createdAfter: string, createdBefore: string) {
	filter.createdAfter = createdAfter;
	filter.createdBefore = createdBefore;
}

export function createJournalPageActions(deps: PageActionDeps) {
	const { client, filter, pagination, reviewPresets, state, filterToUrl, urlToFilter } = deps;

	function setActiveReviewPreset(label: string) {
		state.activeReviewPreset = label;
	}

	function applyFilters() {
		pagination.resetPage();
		filterToUrl();
		void loadJournals();
	}

	function applyManualFilters() {
		setActiveReviewPreset('');
		applyFilters();
	}

	function syncActiveReviewPreset() {
		const matchingPreset = findMatchingReviewPreset(reviewPresets, filter);
		setActiveReviewPreset(matchingPreset?.label ?? '');
	}

	function applyResponse(resp: Awaited<ReturnType<JournalClient['list']>>) {
		const data = resp.data ?? {};
		state.journals = data.journals ?? [];
		pagination.setTotalItems(data.metadata?.total ?? state.journals.length);
		pagination.setPageFromOffset(data.metadata?.offset ?? 0);
	}

	async function loadJournals() {
		state.requestCounter += 1;
		const requestId = state.requestCounter;
		state.loading = true;
		state.errorMessage = '';

		try {
			const resp = await client.list(
				pagination.getOffset(),
				pagination.getPageSize(),
				filter.toQueryParams(),
			);

			if (requestId !== state.requestCounter) {
				return;
			}

			applyResponse(resp);
		} catch {
			if (requestId !== state.requestCounter) {
				return;
			}

			state.errorMessage = 'Unable to load journals. Please try again.';
		} finally {
			if (requestId === state.requestCounter) {
				state.loading = false;
			}
		}
	}

	return {
		applyFilters,
		applyManualFilters,
		syncActiveReviewPreset,
		clearActiveReviewPreset() {
			setActiveReviewPreset('');
		},
		reviewPresetClass(reviewPreset: ReviewPreset) {
			return getReviewPresetClass(state.activeReviewPreset, reviewPreset);
		},
		applyCreatedPreset(preset: CreatedPreset) {
			filter.clear();
			const range = getCreatedPresetRange(preset);
			applyCreatedRange(filter, range.createdAfter, range.createdBefore);
			setActiveReviewPreset('');
			applyFilters();
		},
		applyReviewPreset(reviewPreset: ReviewPreset) {
			filter.clear();
			applyCreatedRange(filter, reviewPreset.createdAfter, reviewPreset.createdBefore);
			filter.reviewed = 'false';
			setActiveReviewPreset(reviewPreset.label);
			applyFilters();
		},
		loadJournals,
		hasError() {
			return state.errorMessage !== '';
		},
		isEmpty() {
			return state.journals.length === 0;
		},
		async prevPage() {
			if (!pagination.hasPrev()) return;
			pagination.prevPage();
			await loadJournals();
		},
		async nextPage() {
			if (!pagination.hasNext()) return;
			pagination.nextPage();
			await loadJournals();
		},
		init() {
			urlToFilter();
			syncActiveReviewPreset();
			void loadJournals();
		},
	};
}

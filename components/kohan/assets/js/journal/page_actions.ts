import type { JournalClient } from '../client/journal';
import type { JournalFilterState } from './filter';
import type { PaginationState } from './pagination';
import type { JournalPageState } from './page_state';

type PageActionDeps = {
	client: JournalClient;
	filter: JournalFilterState;
	pagination: PaginationState;
	state: JournalPageState;
	filterToUrl: () => void;
	urlToFilter: () => void;
	clearActiveReviewPreset: () => void;
	syncActiveReviewPreset: () => void;
};

export function createJournalPageActions(deps: PageActionDeps) {
	const { client, filter, pagination, state, filterToUrl, urlToFilter, clearActiveReviewPreset, syncActiveReviewPreset } = deps;

	function applyFilters() {
		pagination.resetPage();
		filterToUrl();
		void loadJournals();
	}

	function applyManualFilters() {
		clearActiveReviewPreset();
		applyFilters();
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

import type { JournalClient } from '../client/journal';

type PageActionDeps = {
	client: JournalClient;
};

export function createJournalPageActions(deps: PageActionDeps) {
	const { client } = deps;

	async function loadJournals(this: any) {
		this.loading = true;
		this.errorMessage = '';

		try {
			const response = await client.list(
				this.pagination.getOffset(),
				this.pagination.getPageSize(),
				this.filter.toQueryParams(),
			);

			const data = response.data ?? {};
			this.journals = data.journals ?? [];
			this.pagination.setTotalItems(data.metadata?.total ?? this.journals.length);
			this.pagination.setPageFromOffset(data.metadata?.offset ?? 0);
		} finally {
			this.loading = false;
		}
	}

	function applyFilters(this: any) {
		this.pagination.resetPage();
		this.filterToUrl();
		void this.loadJournals();
	}

	function applyManualFilters(this: any) {
		this.clearActiveReviewPreset();
		this.applyFilters();
	}

	return {
		applyFilters,
		applyManualFilters,
		loadJournals,
		hasError(this: any) {
			return this.errorMessage !== '';
		},
		isEmpty(this: any) {
			return this.journals.length === 0;
		},
		async prevPage(this: any) {
			if (!this.pagination.hasPrev()) return;
			this.pagination.prevPage();
			await this.loadJournals();
		},
		async nextPage(this: any) {
			if (!this.pagination.hasNext()) return;
			this.pagination.nextPage();
			await this.loadJournals();
		},
		init(this: any) {
			console.log('[DEBUG] init() called');
			this.urlToFilter();
			this.syncActiveReviewPreset();
			console.log('[DEBUG] init() starting loadJournals');
			void this.loadJournals();
		},
	};
}

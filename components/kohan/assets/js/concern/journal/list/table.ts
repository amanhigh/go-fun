import type { JournalPageProvider, JournalTableConcern } from '../../../types/journal_list_concern';

export function NewTableConcern(pg: JournalPageProvider): JournalTableConcern {
	return {
		journals: [],
		loading: false,
		async loadJournals() {
			this.loading = true;

			try {
				const page = pg();
				const pagination = page.pagination;
				const response = await page.client.list(pagination.getOffset(), pagination.getPageSize(), page.filter);
				const data = response.data ?? {};
				this.journals = data.journals ?? [];
				pagination.setTotalItems(data.metadata?.total ?? this.journals.length);
				pagination.setPageFromOffset(data.metadata?.offset ?? 0);
			} finally {
				this.loading = false;
			}
		},
		isEmpty() { return this.journals.length === 0; },
	};
}

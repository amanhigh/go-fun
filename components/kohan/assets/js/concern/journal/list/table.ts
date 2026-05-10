import { createLoader } from '../../../lib/loader';
import type { Loader } from '../../../lib/loader';
import type { JournalPageProvider, JournalTableConcern } from '../../../types/journal_list_concern';

export function NewTableConcern(pg: JournalPageProvider): JournalTableConcern {
	return {
		journals: [],
		loader: createLoader(),
		async loadJournals() {
			await this.loader.run(async () => {
				const page = pg();
				const pagination = page.pagination;
				const response = await page.client.list(pagination.getOffset(), pagination.getPageSize(), page.filter);
				const data = response.data ?? {};
				this.journals = data.journals ?? [];
				pagination.setTotalItems(data.metadata?.total ?? this.journals.length);
				pagination.setPageFromOffset(data.metadata?.offset ?? 0);
			}, { error: 'Unable to load journals.' });
		},
		isEmpty() { return this.journals.length === 0; },
	};
}

import { createLoader } from '../../../lib/loader';
import type { Loader } from '../../../lib/loader';
import type { JournalPageProvider, JournalTableConcern } from '../../../types/journal_list_concern';

export function NewTableConcern(pg: JournalPageProvider): JournalTableConcern {
	return {
		journals: [],
		loader: createLoader(),
		async loadJournals() {
			const page = pg();
			const pagination = page.pagination;

			const data = await this.loader.loadData(
				() => page.client.list(pagination.getOffset(), pagination.getPageSize(), page.filter),
			);

			if (!data) return;

			this.journals = data.journals;
			pagination.setTotalItems(data.metadata.total);
			pagination.setPageFromOffset(data.metadata.offset);
		},
		isEmpty() { return this.journals.length === 0; },
	};
}

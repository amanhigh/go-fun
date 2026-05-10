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

			const loaded = await this.loader.loadData(
				() => page.client.list(pagination.getOffset(), pagination.getPageSize(), page.filter),
				{ error: 'Unable to load journals.' },
			);

			if (!loaded) return;

			this.journals = loaded.result.journals;
			pagination.setTotalItems(loaded.metadata?.total ?? this.journals.length);
			pagination.setPageFromOffset(loaded.metadata?.offset ?? 0);
		},
		isEmpty() { return this.journals.length === 0; },
	};
}

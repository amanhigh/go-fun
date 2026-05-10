import { createCollection } from '../../../lib/collection';
import { createLoader } from '../../../lib/loader';
import type { Journal } from '../../../types/api/journal/response';
import type { JournalPageProvider } from '../../../types/journal/list';

export function NewTableConcern(pg: JournalPageProvider) {
	return {
		...createCollection<Journal>(),
		loader: createLoader(),

		async load() {
			const page = pg();
			const pagination = page.pagination;

			await this.loader.load(
				() => page.client.list(pagination.getOffset(), pagination.getPageSize(), page.filter),
				(data) => {
					this.sync(data.journals);
					pagination.setTotalItems(data.metadata.total);
					pagination.setPageFromOffset(data.metadata.offset);
				},
			);
		},
	};
}

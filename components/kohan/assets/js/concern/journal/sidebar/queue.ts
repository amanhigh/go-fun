import { createLoader } from '../../../lib/loader';
import { createCollection } from '../../../lib/collection';
import { JournalSortBy, JournalSortOrder } from '../../../types/api/journal/enums';
import { ReviewedFilter } from '../../../types/api/journal/request';
import type { Journal } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function NewReviewQueueConcern(pg: JournalDetailPageProvider) {
	return {
		...createCollection<Journal>(),
		loader: createLoader(),

		async load(this: any) {
			await this.loader.load(
				() => pg().client.list(0, 10, {
					reviewed: ReviewedFilter.PENDING,
					sortBy: JournalSortBy.CREATED_AT,
					sortOrder: JournalSortOrder.ASC,
				}),
				(data: any) => this.sync(data.journals),
			);
		},
	};
}

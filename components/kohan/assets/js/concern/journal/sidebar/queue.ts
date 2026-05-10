import { createLoadableCollectionState } from '../../../lib/collection';
import { JournalSortBy, JournalSortOrder } from '../../../types/api/journal/enums';
import { ReviewedFilter } from '../../../types/api/journal/request';
import type { Journal } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function NewReviewQueueConcern(pg: JournalDetailPageProvider) {
	return {
		...createLoadableCollectionState<Journal>(
			async () => {
				const envelope = await pg().client.list(0, 10, { reviewed: ReviewedFilter.PENDING, sortBy: JournalSortBy.CREATED_AT, sortOrder: JournalSortOrder.ASC });
				return envelope.data.journals;
			},
			'Unable to load review queue.',
		),
	};
}

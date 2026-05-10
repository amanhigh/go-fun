import { createLoadableCollectionState } from '../../../lib/collection';
import { JournalSortBy, JournalSortOrder, ReviewedFilter } from '../../../types/journal_api';
import type { Journal } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

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

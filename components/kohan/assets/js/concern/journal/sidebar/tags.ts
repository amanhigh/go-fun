import { createDeletableSyncedCollectionState } from '../../../lib/collection';
import type { JournalTag } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewTagCollectionConcern(pg: JournalDetailPageProvider) {
	return {
		...createDeletableSyncedCollectionState<JournalTag>(
			() => !!pg().current.journal,
			(tagId) => pg().tagClient.delete(pg().current.journalId, tagId),
		),
		reason() {
			return this.all().filter((tag) => tag.type === 'REASON');
		},
		directional() {
			return this.all().filter((tag) => tag.type === 'DIRECTION');
		},
		management() {
			return this.all().filter((tag) => tag.type === 'MANAGEMENT');
		},
	};
}

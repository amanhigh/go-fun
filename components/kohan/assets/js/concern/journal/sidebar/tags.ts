import { createDeletableSyncedCollectionState } from '../../../lib/collection';
import { JournalTagType } from '../../../types/api/journal/enums';
import type { JournalTag } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function NewTagCollectionConcern(pg: JournalDetailPageProvider) {
	return {
		...createDeletableSyncedCollectionState<JournalTag>(
			() => !!pg().current.journal,
			(tagId) => pg().tagClient.delete(pg().current.journalId, tagId),
		),
		reason() {
			return this.all().filter((tag) => tag.type === JournalTagType.REASON);
		},
		directional() {
			return this.all().filter((tag) => tag.type === JournalTagType.DIRECTION);
		},
		management() {
			return this.all().filter((tag) => tag.type === JournalTagType.MANAGEMENT);
		},
	};
}

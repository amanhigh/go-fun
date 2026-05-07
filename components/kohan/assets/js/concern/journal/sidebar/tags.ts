import { createDeletableSyncedCollectionState } from './collection';
import { normalizeTag } from '../../../lib/tags';
import type { JournalTag } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewTagCollectionConcern(pg: JournalDetailPageProvider) {
	return {
		...createDeletableSyncedCollectionState<JournalTag>(
			() => !!pg().current.journal,
			(tagId) => pg().tagClient.delete(pg().current.journalId, tagId),
		),
		reason() {
			return this.items.filter((tag) => normalizeTag(tag.type ?? '') === 'REASON');
		},
		directional() {
			return this.items.filter((tag) => normalizeTag(tag.type ?? '') === 'DIRECTION');
		},
		management() {
			return this.items.filter((tag) => normalizeTag(tag.type ?? '') === 'MANAGEMENT');
		},
	};
}

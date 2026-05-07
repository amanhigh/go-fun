import { removeById } from '../../../shared/collection';
import { normalizeTag } from '../../../shared/tags';
import type { JournalTag } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewTagCollectionConcern(pg: JournalDetailPageProvider) {
	return {
		items: [] as JournalTag[],
		deletingId: '',

		sync(tags: JournalTag[] | undefined) {
			this.items = [...(tags ?? [])];
		},
		all() {
			return this.items;
		},
		reason() {
			return this.items.filter((tag) => normalizeTag(tag.type ?? '') === 'REASON');
		},
		directional() {
			return this.items.filter((tag) => normalizeTag(tag.type ?? '') === 'DIRECTION');
		},
		management() {
			return this.items.filter((tag) => normalizeTag(tag.type ?? '') === 'MANAGEMENT');
		},
		hasTags() {
			return this.all().length > 0;
		},
		async delete(tagId: string) {
			if (!pg().current.journal) return;
			this.deletingId = tagId;
			try {
				await pg().tagClient.delete(pg().current.journalId, tagId);
				this.items = removeById(this.items, tagId);

			} finally {
				this.deletingId = '';
			}
		},
	};
}

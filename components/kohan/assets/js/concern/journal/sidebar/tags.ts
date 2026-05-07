import { removeById } from '../../../shared/collection';
import { getErrorMessage } from '../../../shared/error';
import { normalizeTag } from '../../../shared/tags';
import type { JournalTag } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

// Tag grouping rules — defines which normalized tag types belong to each header section.
const REASON_TAG_TYPES = new Set(['REASON', 'MANAGEMENT']);
const DIRECTIONAL_TAG_TYPES = new Set(['DIRECTION', 'LEGACY']);
const MANAGEMENT_TAG_TYPES = new Set(['MANAGEMENT']);

function hasTagType(tag: JournalTag, types: Set<string>): boolean {
	return types.has(normalizeTag(tag.type ?? ''));
}

export function NewTagCollectionConcern(pg: JournalDetailPageProvider) {
	return {
		items: [] as JournalTag[],
		deletingId: '',
		deleteError: '',

		sync(tags: JournalTag[] | undefined) {
			this.items = [...(tags ?? [])];
		},
		all() {
			return this.items ?? [];
		},
		reason() {
			return (this.items ?? []).filter((tag: JournalTag) => hasTagType(tag, REASON_TAG_TYPES));
		},
		directional() {
			return (this.items ?? []).filter((tag: JournalTag) => hasTagType(tag, DIRECTIONAL_TAG_TYPES));
		},
		management() {
			return (this.items ?? []).filter((tag: JournalTag) => hasTagType(tag, MANAGEMENT_TAG_TYPES));
		},
		hasTags() {
			return this.all().length > 0;
		},
		async delete(tagId: string) {
			if (!pg().current.journal || this.deletingId) return;
			this.deletingId = tagId;
			this.deleteError = '';
			try {
				await pg().tagClient.delete(pg().current.journalId, tagId);
				this.items = removeById(this.items ?? [], tagId);
			} catch (err) {
				this.deleteError = getErrorMessage(err, 'Unable to delete tag.');
			} finally {
				this.deletingId = '';
			}
		},
	};
}

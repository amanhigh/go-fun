import { removeById } from '../../../shared/collection';
import { getErrorMessage } from '../../../shared/error';
import { normalizeTag } from '../../../shared/tags';
import type { JournalTag } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewTagCollectionConcern(pg: JournalDetailPageProvider) {
	return {
		items: [] as JournalTag[],
		deletingId: '',

		sync(this: any, tags: JournalTag[] | undefined) {
			this.items = [...(tags ?? [])];
		},
		all(this: any) {
			return this.items ?? [];
		},
		reason(this: any) {
			return (this.items ?? []).filter(
				(tag: JournalTag) => normalizeTag(tag.type ?? '') === 'REASON' || normalizeTag(tag.type ?? '') === 'MANAGEMENT',
			);
		},
		directional(this: any) {
			return (this.items ?? []).filter(
				(tag: JournalTag) => normalizeTag(tag.type ?? '') === 'DIRECTION' || normalizeTag(tag.type ?? '') === 'LEGACY',
			);
		},
		reasonLabel(this: any, tag: JournalTag) {
			const name = tag.tag ?? '';
			const prefix = name.toLowerCase().includes('trend') ? '📈 ' : '⚡ ';
			const override = tag.override ? ` → ${tag.override}` : '';
			return `${prefix}${name}${override}`;
		},
		directionalLabel(this: any, tag: JournalTag) {
			return tag.tag ?? '';
		},
		management(this: any) {
			return (this.items ?? []).filter((tag: JournalTag) => normalizeTag(tag.type ?? '') === 'MANAGEMENT');
		},
		async delete(this: any, tagId: string) {
			if (!pg().current.journal || this.deletingId) return;
			this.deletingId = tagId;
			try {
				await pg().tagClient.delete(pg().current.journalId, tagId);
				this.items = removeById(this.items ?? [], tagId);
			} catch (err) {
				getErrorMessage(err, 'Unable to delete tag.');
			} finally {
				this.deletingId = '';
			}
		},
	};
}

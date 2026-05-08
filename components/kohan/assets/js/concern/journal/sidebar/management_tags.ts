import { createAsyncFeedback } from '../../../lib/async_feedback';
import { normalizeTag } from '../../../lib/tags';
import { managementTagPresets, managementTagTone } from '../../../lib/management_tags';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewManagementTagsConcern(pg: JournalDetailPageProvider) {
	return {
		...createAsyncFeedback(),
		presets: managementTagPresets,
		pendingValue: '',

		hasBar() {
			return normalizeTag(pg().current.journal?.type ?? '') === 'TAKEN';
		},
		hasTag(value: string) {
			const normalizedValue = normalizeTag(value);
			return pg().sidebar.tags.management().some((tag: JournalTag) => normalizeTag(tag.tag ?? '') === normalizedValue);
		},
		buttonClass(value: string) {
			const tagKey = normalizeTag(value);
			const tone = managementTagTone(value);
			const isActive = this.hasTag(value);
			const isPending = this.submitting && normalizeTag(this.pendingValue) === tagKey;
			const baseClass = isActive ? `journal-management-active-${tone}` : `journal-management-base-${tone}`;
			return isPending ? `journal-management-pending ${baseClass}` : baseClass;
		},
		async submit(tagValue: string) {
			if (!pg().current.journal || this.submitting) return;
			this.pendingValue = tagValue;
			await this.run(
				() => this.addTag(tagValue),
				`${normalizeTag(tagValue)} tag added.`,
				'Unable to save management tag.',
			);
			this.pendingValue = '';
		},

		async addTag(tagValue: string) {
			const payload: JournalTagRequest = { tag: tagValue, type: 'MANAGEMENT' };
			const page = pg();
			const envelope = await page.tagClient.create(page.current.journalId, payload);
			page.sidebar.tags.prepend(envelope.data as JournalTag);
		},
	};
}

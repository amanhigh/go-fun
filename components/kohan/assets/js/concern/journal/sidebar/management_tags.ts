import { createAsyncFeedbackState } from '../../../lib/async_feedback';
import { getErrorMessage } from '../../../lib/error';
import { normalizeTag } from '../../../lib/tags';
import { managementTagPresets, managementTagTone } from '../../../lib/management_tags';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewManagementTagsConcern(pg: JournalDetailPageProvider) {
	return {
		...createAsyncFeedbackState(),
		presets: managementTagPresets,
		pendingValue: '',

		get feedbackClass(): string {
			return this.messageType === 'success' ? 'journal-feedback-success' : 'journal-feedback-error';
		},

		hasBar() {
			return normalizeTag(pg().current.journal?.type ?? '') === 'TAKEN';
		},
		hasTag(value: string) {
			const normalizedValue = normalizeTag(value);
			return pg().sidebar.tags.management().some((tag: JournalTag) => normalizeTag(tag.tag ?? '') === normalizedValue);
		},
		buttonClass(value: string) {
			const tagKey = normalizeTag(value);
			const tone = managementTagTone(tagKey);
			const isActive = this.hasTag(value);
			const isPending = this.submitting && normalizeTag(this.pendingValue) === tagKey;
			const baseClass = isActive ? `journal-management-active-${tone}` : `journal-management-base-${tone}`;
			return isPending ? `journal-management-pending ${baseClass}` : baseClass;
		},
		async submit(tagValue: string) {
			if (!pg().current.journal || this.submitting) return;
			this.submitting = true;
			this.pendingValue = tagValue;
			this.message = '';
			this.messageType = 'error';
			try {
				const payload: JournalTagRequest = {
					tag: tagValue,
					type: 'MANAGEMENT',
				};
				const envelope = await pg().tagClient.create(pg().current.journalId, payload);
				pg().sidebar.tags.prepend(envelope.data as JournalTag);
				this.messageType = 'success';
				this.message = `${normalizeTag(tagValue)} tag added.`;
			} catch (err) {
				this.message = getErrorMessage(err, 'Unable to save management tag.');
				this.messageType = 'error';
			} finally {
				this.submitting = false;
				this.pendingValue = '';
			}
		},
	};
}

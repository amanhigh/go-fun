import { createAsyncFeedbackState } from '../../../lib/async_feedback';
import { getErrorMessage } from '../../../lib/error';
import { normalizeTag } from '../../../lib/tags';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

const managementTagPresets = [
	{ value: 'ntr', label: 'NTR', tone: 'emerald' },
	{ value: 'enl', label: 'ENL', tone: 'sky' },
	{ value: 'slt', label: 'SLT', tone: 'rose' },
	{ value: 'fz', label: 'FZ', tone: 'violet' },
	{ value: 'nbe', label: 'NBE', tone: 'amber' },
	{ value: 'ws', label: 'WS', tone: 'slate' },
	{ value: 'important', label: 'IMPORTANT', tone: 'fuchsia' },
	{ value: 'be', label: 'BE', tone: 'orange' },
] as const;

function toneForValue(value: string): string {
	const preset = managementTagPresets.find((p) => normalizeTag(p.value) === normalizeTag(value));
	return preset?.tone ?? 'slate';
}

export function NewManagementTagsConcern(pg: JournalDetailPageProvider) {
	return {
		...createAsyncFeedbackState('submitting', 'message', 'messageType'),
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
			const normalizedValue = normalizeTag(value);
			const tone = toneForValue(normalizedValue);
			const isActive = this.hasTag(value);
			const isPending = this.submitting && normalizeTag(this.pendingValue) === normalizedValue;
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

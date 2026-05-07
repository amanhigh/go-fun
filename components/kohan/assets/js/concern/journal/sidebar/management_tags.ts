import { createAsyncFeedbackState } from '../../../shared/async_feedback';
import { getErrorMessage } from '../../../shared/error';
import { normalizeTag } from '../../../shared/tags';
import { prependById } from '../../../shared/collection';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

const managementTagPresets = [
	{ value: 'ntr', label: 'NTR' },
	{ value: 'enl', label: 'ENL' },
	{ value: 'slt', label: 'SLT' },
	{ value: 'fz', label: 'FZ' },
	{ value: 'nbe', label: 'NBE' },
	{ value: 'ws', label: 'WS' },
	{ value: 'important', label: 'IMPORTANT' },
	{ value: 'be', label: 'BE' },
] as const;

const managementTagToneMap: Record<string, string> = {
	NTR: 'emerald',
	ENL: 'sky',
	SLT: 'rose',
	FZ: 'violet',
	NBE: 'amber',
	WS: 'slate',
	IMPORTANT: 'fuchsia',
	BE: 'orange',
};

export function createManagementTagsState() {
	return {
		...createAsyncFeedbackState('submitting', 'message', 'messageType'),
		presets: managementTagPresets,
		pendingValue: '',
	};
}

export function NewManagementTagsConcern(pg: JournalDetailPageProvider) {
	return {
		...createManagementTagsState(),

		get feedbackClass(): string {
			return this.messageType === 'success' ? 'journal-feedback-success' : 'journal-feedback-error';
		},

		hasBar(this: any) {
			return normalizeTag(pg().current.journal?.type ?? '') === 'TAKEN';
		},
		hasTag(this: any, value: string) {
			const normalizedValue = normalizeTag(value);
			return pg().sidebar.tags.management().some((tag: JournalTag) => normalizeTag(tag.tag ?? '') === normalizedValue);
		},
		buttonClass(this: any, value: string) {
			const normalizedValue = normalizeTag(value);
			const tone = managementTagToneMap[normalizedValue] ?? 'slate';
			const isActive = this.hasTag(value);
			const isPending = this.submitting && normalizeTag(this.pendingValue) === normalizedValue;
			const baseClass = isActive ? `journal-management-active-${tone}` : `journal-management-base-${tone}`;
			return isPending ? `journal-management-pending ${baseClass}` : baseClass;
		},
		async submit(this: any, tagValue: string) {
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
				pg().sidebar.tags.items = prependById(pg().sidebar.tags.items ?? [], envelope.data as JournalTag);
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

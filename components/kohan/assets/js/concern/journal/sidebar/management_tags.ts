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
			return this.messageType === 'success' ? 'text-emerald-700' : 'text-rose-700';
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
			if (isActive) {
				return {
					emerald: 'border-emerald-400 bg-emerald-100 text-emerald-900 opacity-90',
					sky: 'border-sky-400 bg-sky-100 text-sky-900 opacity-90',
					rose: 'border-rose-400 bg-rose-100 text-rose-900 opacity-90',
					violet: 'border-violet-400 bg-violet-100 text-violet-900 opacity-90',
					amber: 'border-amber-400 bg-amber-100 text-amber-900 opacity-90',
					slate: 'border-slate-400 bg-slate-100 text-slate-900 opacity-90',
					fuchsia: 'border-fuchsia-400 bg-fuchsia-100 text-fuchsia-900 opacity-90',
					orange: 'border-orange-400 bg-orange-100 text-orange-900 opacity-90',
				}[tone] ?? 'border-slate-400 bg-slate-100 text-slate-900 opacity-90';
			}

			const baseClass = {
				emerald: 'border-emerald-300 bg-emerald-50 text-emerald-800 hover:bg-emerald-100 focus:border-emerald-400 focus:ring-emerald-200',
				sky: 'border-sky-300 bg-sky-50 text-sky-800 hover:bg-sky-100 focus:border-sky-400 focus:ring-sky-200',
				rose: 'border-rose-300 bg-rose-50 text-rose-800 hover:bg-rose-100 focus:border-rose-400 focus:ring-rose-200',
				violet: 'border-violet-300 bg-violet-50 text-violet-800 hover:bg-violet-100 focus:border-violet-400 focus:ring-violet-200',
				amber: 'border-amber-300 bg-amber-50 text-amber-800 hover:bg-amber-100 focus:border-amber-400 focus:ring-amber-200',
				slate: 'border-slate-300 bg-slate-50 text-slate-800 hover:bg-slate-100 focus:border-slate-400 focus:ring-slate-200',
				fuchsia: 'border-fuchsia-300 bg-fuchsia-50 text-fuchsia-800 hover:bg-fuchsia-100 focus:border-fuchsia-400 focus:ring-fuchsia-200',
				orange: 'border-orange-300 bg-orange-50 text-orange-800 hover:bg-orange-100 focus:border-orange-400 focus:ring-orange-200',
			}[tone] ?? 'border-slate-300 bg-slate-50 text-slate-800 hover:bg-slate-100 focus:border-slate-400 focus:ring-slate-200';

			return isPending ? `opacity-70 ${baseClass}` : baseClass;
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

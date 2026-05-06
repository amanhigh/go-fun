import { createAsyncFeedbackState } from '../../../shared/async_feedback';
import { prependById, removeById } from '../../../shared/collection';
import { getErrorMessage } from '../../../shared/error';
import { normalizeTag } from '../../../shared/tags';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider, TagsState } from '../../../types/journal_detail_concern';

export const managementTagPresets = [
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

export function createTagsState(tagPresets: typeof managementTagPresets): TagsState {
	return {
		...createAsyncFeedbackState('managementTagSubmitting', 'managementTagMessage', 'managementTagMessageType'),
		...createAsyncFeedbackState('reasonTagSubmitting', 'reasonTagMessage', 'reasonTagMessageType'),
		managementTagPresets: tagPresets,
		managementTagPendingValue: '',
		reasonTagInput: '',
		reasonTagOverride: '',
		tagItems: [],
		tagDeletingId: '',
	};
}

export function NewTagsConcern(pg: JournalDetailPageProvider) {
	return {
		syncTags(this: any, tags: JournalTag[] | undefined) {
			this.tagItems = [...(tags ?? [])];
		},
		reasonTags(this: any) {
			return (this.tagItems ?? []).filter((tag: JournalTag) => normalizeTag(tag.type ?? '') === 'REASON' || normalizeTag(tag.type ?? '') === 'MANAGEMENT');
		},
		tags(this: any) {
			return this.tagItems ?? [];
		},
		directionalTags(this: any) {
			return (this.tagItems ?? []).filter((tag: JournalTag) => normalizeTag(tag.type ?? '') === 'DIRECTION' || normalizeTag(tag.type ?? '') === 'LEGACY');
		},
		reasonTagLabel(this: any, tag: JournalTag) {
			const name = tag.tag ?? '';
			const prefix = name.toLowerCase().includes('trend') ? '📈 ' : '⚡ ';
			const override = tag.override ? ` → ${tag.override}` : '';
			return `${prefix}${name}${override}`;
		},
		directionalTagLabel(this: any, tag: JournalTag) {
			return tag.tag ?? '';
		},
		managementTags(this: any) {
			return (this.tagItems ?? []).filter((tag: JournalTag) => normalizeTag(tag.type ?? '') === 'MANAGEMENT');
		},
		hasManagementBar(this: any) {
			return normalizeTag(pg().journal?.type ?? '') === 'TAKEN';
		},
		hasManagementTag(this: any, value: string) {
			const normalizedValue = normalizeTag(value);
			return this.managementTags().some((tag: JournalTag) => normalizeTag(tag.tag ?? '') === normalizedValue);
		},
		managementTagButtonClass(this: any, value: string) {
			const normalizedValue = normalizeTag(value);
			const tone = managementTagToneMap[normalizedValue] ?? 'slate';
			const isActive = this.hasManagementTag(value);
			const isPending = this.managementTagSubmitting && normalizeTag(this.managementTagPendingValue) === normalizedValue;
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
		focusReasonTagOverride(this: any) {
			pg().$nextTick(() => {
				pg().$refs?.reasonTagOverride?.focus?.();
			});
		},
		async submitReasonTag(this: any) {
			if (!pg().journal || this.reasonTagSubmitting) return;
			const tag = this.reasonTagInput.trim();
			if (!tag) {
				this.reasonTagMessage = 'Tag is required.';
				this.reasonTagMessageType = 'error';
				return;
			}
			const override = this.reasonTagOverride.trim();
			this.reasonTagSubmitting = true;
			this.reasonTagMessage = '';
			this.reasonTagMessageType = 'error';
			try {
				const payload: JournalTagRequest = {
					tag,
					type: 'REASON',
					...(override ? { override } : {}),
				};
				const envelope = await pg().tagClient.create(pg().journalId, payload);
				this.tagItems = prependById(this.tagItems ?? [], envelope.data);
				this.reasonTagInput = '';
				this.reasonTagOverride = '';
				this.reasonTagMessageType = 'success';
				this.reasonTagMessage = 'Reason tag added.';
			} catch (err) {
				this.reasonTagMessage = getErrorMessage(err, 'Unable to save reason tag.');
				this.reasonTagMessageType = 'error';
			} finally {
				this.reasonTagSubmitting = false;
			}
		},
		async submitManagementTag(this: any, tagValue: string) {
			if (!pg().journal || this.managementTagSubmitting || !this.hasManagementBar() || !this.hasManagementTag || this.hasManagementTag(tagValue)) return;
			this.managementTagSubmitting = true;
			this.managementTagPendingValue = tagValue;
			this.managementTagMessage = '';
			this.managementTagMessageType = 'error';
			try {
				const payload: JournalTagRequest = {
					tag: tagValue,
					type: 'MANAGEMENT',
				};
				const envelope = await pg().tagClient.create(pg().journalId, payload);
				this.tagItems = prependById(this.tagItems ?? [], envelope.data);
				this.managementTagMessageType = 'success';
				this.managementTagMessage = `${normalizeTag(tagValue)} tag added.`;
			} catch (err) {
				this.managementTagMessage = getErrorMessage(err, 'Unable to save management tag.');
				this.managementTagMessageType = 'error';
			} finally {
				this.managementTagSubmitting = false;
				this.managementTagPendingValue = '';
			}
		},
		async deleteTag(this: any, tagId: string) {
			if (!pg().journal || this.tagDeletingId) return;
			this.tagDeletingId = tagId;
			this.reasonTagMessage = '';
			this.reasonTagMessageType = 'error';
			try {
				await pg().tagClient.delete(pg().journalId, tagId);
				this.tagItems = removeById(this.tagItems ?? [], tagId);
				this.reasonTagMessageType = 'success';
				this.reasonTagMessage = 'Tag deleted.';
			} catch (err) {
				this.reasonTagMessage = getErrorMessage(err, 'Unable to delete tag.');
				this.reasonTagMessageType = 'error';
			} finally {
				this.tagDeletingId = '';
			}
		},
	};
}

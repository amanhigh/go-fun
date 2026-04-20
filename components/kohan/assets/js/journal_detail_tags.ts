import type { Envelope, JournalTag, JournalTagCreate } from './journal_models';
import { normalizeTag } from './journal_detail_formatters';

export const managementTagPresets = [
	{ value: 'ntr', label: 'NTR' },
	{ value: 'enl', label: 'ENL' },
	{ value: 'slt', label: 'SLT' },
	{ value: 'fz', label: 'FZ' },
	{ value: 'nbe', label: 'NBE' },
	{ value: 'ws', label: 'WS' },
	{ value: 'important', label: 'IMPORTANT' },
	{ value: 'be', label: 'BE' },
];

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

export function createJournalDetailTags() {
	return {
		reasonTags(this: any) {
			return (this.journal?.tags ?? []).filter((tag: JournalTag) => normalizeTag(tag.type ?? '') === 'REASON');
		},
		deletableTags(this: any) {
			return this.journal?.tags ?? [];
		},
		managementTags(this: any) {
			return (this.journal?.tags ?? []).filter((tag: JournalTag) => normalizeTag(tag.type ?? '') === 'MANAGEMENT');
		},
		hasManagementBar(this: any) {
			return normalizeTag(this.journal?.type ?? '') === 'TAKEN';
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
			this.$nextTick(() => {
				this.$refs?.reasonTagOverride?.focus?.();
			});
		},
		async submitReasonTag(this: any) {
			if (!this.journal || this.reasonTagSubmitting) return;
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
				const payload: JournalTagCreate = {
					tag,
					type: 'REASON',
					...(override ? { override } : {}),
				};
				const response = await fetch(`/v1/api/journals/${this.journalId}/tags`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(payload),
				});
				if (!response.ok) throw new Error(response.status === 404 ? 'Journal not found' : 'Failed to save reason tag');
				const envelope = (await response.json()) as Envelope<JournalTag>;
				const tags = this.journal.tags ?? [];
				this.journal.tags = [envelope.data, ...tags.filter((item: JournalTag) => item.id !== envelope.data.id)];
				this.reasonTagInput = '';
				this.reasonTagOverride = '';
				this.reasonTagMessageType = 'success';
				this.reasonTagMessage = 'Reason tag added.';
			} catch (err) {
				this.reasonTagMessage = err instanceof Error ? err.message : 'Unable to save reason tag.';
				this.reasonTagMessageType = 'error';
			} finally {
				this.reasonTagSubmitting = false;
			}
		},
		async submitManagementTag(this: any, tagValue: string) {
			if (!this.journal || this.managementTagSubmitting || !this.hasManagementBar() || this.hasManagementTag(tagValue)) return;
			this.managementTagSubmitting = true;
			this.managementTagPendingValue = tagValue;
			this.managementTagMessage = '';
			this.managementTagMessageType = 'error';
			try {
				const payload: JournalTagCreate = {
					tag: tagValue,
					type: 'MANAGEMENT',
				};
				const response = await fetch(`/v1/api/journals/${this.journalId}/tags`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(payload),
				});
				if (!response.ok) throw new Error(response.status === 404 ? 'Journal not found' : 'Failed to save management tag');
				const envelope = (await response.json()) as Envelope<JournalTag>;
				const tags = this.journal.tags ?? [];
				this.journal.tags = [envelope.data, ...tags.filter((item: JournalTag) => item.id !== envelope.data.id)];
				this.managementTagMessageType = 'success';
				this.managementTagMessage = `${normalizeTag(tagValue)} tag added.`;
			} catch (err) {
				this.managementTagMessage = err instanceof Error ? err.message : 'Unable to save management tag.';
				this.managementTagMessageType = 'error';
			} finally {
				this.managementTagSubmitting = false;
				this.managementTagPendingValue = '';
			}
		},
		async deleteTag(this: any, tagId: string) {
			if (!this.journal || this.tagDeletingId) return;
			this.tagDeletingId = tagId;
			this.reasonTagMessage = '';
			this.reasonTagMessageType = 'error';
			try {
				const response = await fetch(`/v1/api/journals/${this.journalId}/tags/${tagId}`, {
					method: 'DELETE',
				});
				if (!response.ok) throw new Error(response.status === 404 ? 'Tag not found' : 'Failed to delete tag');
				this.journal.tags = (this.journal.tags ?? []).filter((tag: JournalTag) => tag.id !== tagId);
				this.reasonTagMessageType = 'success';
				this.reasonTagMessage = 'Tag deleted.';
			} catch (err) {
				this.reasonTagMessage = err instanceof Error ? err.message : 'Unable to delete tag.';
				this.reasonTagMessageType = 'error';
			} finally {
				this.tagDeletingId = '';
			}
		},
	};
}

import {
	type Journal,
	type JournalList,
	type Envelope,
	type JournalNote,
	type JournalNoteCreate,
	type JournalTag,
	type JournalTagCreate,
	type JournalReviewStatusResponse,
	type JournalReviewUpdate,
} from './journal_models';
import { createImageHelper } from './journal_images';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalDetailPage>): void;
};

const defaultBadgeClass = 'border-slate-300 bg-slate-50 text-slate-700';

const badgeClassMap: Record<string, Record<string, string>> = {
	status: {
		SUCCESS: 'border-emerald-300 bg-emerald-50 text-emerald-800',
		FAIL: 'border-rose-300 bg-rose-50 text-rose-800',
		RUNNING: 'border-sky-300 bg-sky-50 text-sky-800',
		SET: 'border-amber-300 bg-amber-50 text-amber-800',
		JUST_LOSS: 'border-rose-300 bg-rose-50 text-rose-800',
		BROKEN: 'border-violet-300 bg-violet-50 text-violet-800',
		MISSED: 'border-slate-300 bg-slate-50 text-slate-700',
		REJECTED: 'border-violet-300 bg-violet-50 text-violet-800',
	},
	type: {
		REJECTED: 'border-rose-300 bg-rose-50 text-rose-800',
		RESULT: 'border-emerald-300 bg-emerald-50 text-emerald-800',
		SET: 'border-indigo-300 bg-indigo-50 text-indigo-800',
	},
};

const normalizeTag = (value: string): string => (value ?? '').trim().toUpperCase();
const shortMonthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
const managementTagPresets = [
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

function journalDetailPage() {
	const image = createImageHelper();
	return {
		journalId: '',
		journal: null as Journal | null,
		selectedImageIndex: -1,
		loading: false,
		errorMessage: '',
		reviewSubmitting: false,
		noteSubmitting: false,
		noteDeletingId: '' as string,
		noteContent: '',
		reviewMessage: '',
		reviewMessageType: 'error',
		noteMessage: '',
		noteMessageType: 'error',
		reviewQueue: [] as Journal[],
		reviewQueueLoading: false,
		reviewQueueError: '',
		managementTagPresets,
		managementTagSubmitting: false,
		managementTagPendingValue: '',
		managementTagMessage: '',
		managementTagMessageType: 'error',
		reasonTagInput: '',
		reasonTagOverride: '',
		reasonTagSubmitting: false,
		tagDeletingId: '',
		reasonTagMessage: '',
		reasonTagMessageType: 'error',
		normalizeStatus: normalizeTag,
		statusBadgeClass: (value: string) => badgeClassMap.status[normalizeTag(value)] ?? defaultBadgeClass,
		typeBadgeClass: (value: string) => badgeClassMap.type[normalizeTag(value)] ?? defaultBadgeClass,
		feedbackClass: (type: string) =>
			type === 'success' ? 'text-emerald-700' : 'text-rose-700',
		reviewQueueItemClass: (value: string) => {
			const journalType = normalizeTag(value);
			if (journalType === 'TAKEN') {
				return 'border-emerald-300 bg-emerald-50/70 hover:bg-emerald-100/80 text-emerald-900';
			}
			if (journalType === 'REJECTED') {
				return 'border-rose-300 bg-rose-50/70 hover:bg-rose-100/80 text-rose-900';
			}
			return 'border-border bg-muted/30 hover:bg-muted/70 hover:text-foreground';
		},
		reviewToggleLabel(this: any) {
			return this.journal?.reviewed_at ? 'Mark Pending' : 'Mark Reviewed';
		},
		reviewButtonClass(this: any) {
			return this.journal?.reviewed_at
				? 'border-amber-300 bg-amber-50 text-amber-800 hover:bg-amber-100 focus:border-amber-400 focus:ring-amber-200'
				: 'border-emerald-300 bg-emerald-50 text-emerald-800 hover:bg-emerald-100 focus:border-emerald-400 focus:ring-emerald-200';
		},
		quickReviewStatus(this: any) {
			const journalType = normalizeTag(this.journal?.type ?? '');
			if (journalType === 'TAKEN') return 'JUST_LOSS';
			if (journalType === 'REJECTED') return 'BROKEN';
			return '';
		},
		quickReviewLabel(this: any) {
			const status = this.quickReviewStatus();
			if (status === 'JUST_LOSS') return 'Mark Just Loss';
			if (status === 'BROKEN') return 'Mark Broken';
			return 'Update Status';
		},
		hasQuickReviewAction(this: any) {
			const targetStatus = this.quickReviewStatus();
			if (!targetStatus || !this.journal) return false;
			return normalizeTag(this.journal.status) !== targetStatus;
		},
		quickReviewButtonClass(this: any) {
			return this.quickReviewStatus() === 'JUST_LOSS'
				? 'border-rose-300 bg-rose-50 text-rose-800 hover:bg-rose-100 focus:border-rose-400 focus:ring-rose-200'
				: 'border-violet-300 bg-violet-50 text-violet-800 hover:bg-violet-100 focus:border-violet-400 focus:ring-violet-200';
		},
		applyReviewUpdate(this: any, payload: JournalReviewUpdate, successMessage: string, errorMessage: string) {
			return (async () => {
				const response = await fetch(`/v1/api/journals/${this.journalId}`, {
					method: 'PATCH',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(payload),
				});
				if (!response.ok) throw new Error(response.status === 404 ? 'Journal not found' : errorMessage);
				const envelope = (await response.json()) as Envelope<JournalReviewStatusResponse>;
				if (this.journal) {
					this.journal.status = envelope.data.status;
					this.journal.reviewed_at = envelope.data.reviewed_at;
				}
				this.reviewMessageType = 'success';
				this.reviewMessage = successMessage;
				await this.loadReviewQueue();
			})();
		},
		// Image helpers
		timeframeChipClass: image.chipClass,
		sortedImages(this: any) {
			return image.sorted(this.journal?.images);
		},
		resolveImageSrc: image.resolve,
		previewImageSrc(this: any) {
			return image.resolve(this.previewImage()?.file_name ?? '', this.previewImage()?.created_at);
		},
		previewImageLabel(this: any) {
			return image.label(this.previewImage());
		},
		previewImageCounter(this: any) {
			return image.counter(this.selectedImageIndex, this.sortedImages().length);
		},
		sortedNotes(this: any) {
			return [...(this.journal?.notes ?? [])].sort((left: JournalNote, right: JournalNote) => {
				const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
				const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
				return rightTime - leftTime;
			});
		},
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
		async loadJournal() {
			this.loading = true;
			this.errorMessage = '';
			try {
				const response = await fetch(`/v1/api/journals/${this.journalId}`);
				if (!response.ok) throw new Error(response.status === 404 ? 'Journal not found' : 'Failed to load journal');
				const envelope = (await response.json()) as Envelope<Journal>;
				this.journal = envelope.data ?? null;
			} catch (err) {
				this.errorMessage = err instanceof Error ? err.message : 'Unable to load journal details. Please try again.';
			} finally {
				this.loading = false;
			}
		},
		localToday() {
			const today = new Date();
			const year = today.getFullYear();
			const month = `${today.getMonth() + 1}`.padStart(2, '0');
			const day = `${today.getDate()}`.padStart(2, '0');
			return `${year}-${month}-${day}`;
		},
		async toggleReview() {
			if (!this.journal || this.reviewSubmitting) return;
			this.reviewSubmitting = true;
			this.reviewMessage = '';
			this.reviewMessageType = 'error';
			try {
				const reviewedAt = this.journal.reviewed_at ? null : this.localToday();
				const payload: JournalReviewUpdate = { reviewed_at: reviewedAt };
				await this.applyReviewUpdate(
					payload,
					reviewedAt ? 'Journal marked reviewed.' : 'Journal marked not reviewed.',
					'Failed to update review date',
				);
			} catch (err) {
				this.reviewMessage = err instanceof Error ? err.message : 'Unable to update review date.';
				this.reviewMessageType = 'error';
			} finally {
				this.reviewSubmitting = false;
			}
		},
		async applyQuickReviewStatus() {
			if (!this.journal || this.reviewSubmitting || !this.hasQuickReviewAction()) return;
			const status = this.quickReviewStatus();
			if (!status) return;
			this.reviewSubmitting = true;
			this.reviewMessage = '';
			this.reviewMessageType = 'error';
			try {
				await this.applyReviewUpdate(
					{ status, reviewed_at: this.localToday() },
					`${this.quickReviewLabel()} applied and journal marked reviewed.`,
					'Failed to update journal status',
				);
			} catch (err) {
				this.reviewMessage = err instanceof Error ? err.message : 'Unable to update journal status.';
				this.reviewMessageType = 'error';
			} finally {
				this.reviewSubmitting = false;
			}
		},
		async submitNote() {
			if (!this.journal || this.noteSubmitting) return;
			const content = this.noteContent.trim();
			if (!content) {
				this.noteMessage = 'Note content is required.';
				this.noteMessageType = 'error';
				return;
			}
			this.noteSubmitting = true;
			this.noteMessage = '';
			this.noteMessageType = 'error';
			try {
				const payload: JournalNoteCreate = {
					status: this.journal.status,
					content,
					format: 'MARKDOWN',
				};
				const response = await fetch(`/v1/api/journals/${this.journalId}/notes`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(payload),
				});
				if (!response.ok) throw new Error(response.status === 404 ? 'Journal not found' : 'Failed to save note');
				const envelope = (await response.json()) as Envelope<JournalNote>;
				const notes = this.journal.notes ?? [];
				notes.unshift(envelope.data);
				this.journal.notes = notes;
				this.noteContent = '';
				this.noteMessageType = 'success';
				this.noteMessage = 'Note added.';
			} catch (err) {
				this.noteMessage = err instanceof Error ? err.message : 'Unable to save note.';
				this.noteMessageType = 'error';
			} finally {
				this.noteSubmitting = false;
			}
		},
		async submitReasonTag() {
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
		async submitManagementTag(tagValue: string) {
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
		async deleteTag(tagId: string) {
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
		async deleteNote(noteId: string) {
			if (!this.journal || this.noteDeletingId) return;
			this.noteDeletingId = noteId;
			this.noteMessage = '';
			this.noteMessageType = 'error';
			try {
				const response = await fetch(`/v1/api/journals/${this.journalId}/notes/${noteId}`, {
					method: 'DELETE',
				});
				if (!response.ok) throw new Error(response.status === 404 ? 'Note not found' : 'Failed to delete note');
				this.journal.notes = (this.journal.notes ?? []).filter((note) => note.id !== noteId);
				this.noteMessageType = 'success';
				this.noteMessage = 'Note deleted.';
			} catch (err) {
				this.noteMessage = err instanceof Error ? err.message : 'Unable to delete note.';
				this.noteMessageType = 'error';
			} finally {
				this.noteDeletingId = '';
			}
		},
		hasError(this: any) {
			return this.errorMessage !== '';
		},
		hasImagePreview(this: any) {
			return this.selectedImageIndex >= 0 && this.selectedImageIndex < this.sortedImages().length;
		},
		openImagePreview(this: any, index: number) {
			this.selectedImageIndex = index;
		},
		closeImagePreview(this: any) {
			this.selectedImageIndex = -1;
		},
		canPrevImage(this: any) {
			return this.selectedImageIndex > 0;
		},
		canNextImage(this: any) {
			return this.selectedImageIndex >= 0 && this.selectedImageIndex < this.sortedImages().length - 1;
		},
		prevImage(this: any, wrap = false) {
			if (this.canPrevImage()) this.selectedImageIndex--;
			else if (wrap && this.sortedImages().length > 0) this.selectedImageIndex = this.sortedImages().length - 1;
		},
		nextImage(this: any, wrap = false) {
			if (this.canNextImage()) this.selectedImageIndex++;
			else if (wrap && this.sortedImages().length > 0) this.selectedImageIndex = 0;
		},
		previewImage(this: any) {
			return this.sortedImages()[this.selectedImageIndex] ?? null;
		},
		previewImageTimeframe(this: any) {
			return this.previewImage()?.timeframe ?? '';
		},
		formatTimestamp: (value: string | null | undefined) => {
			if (!value) return '—';
			const parsed = new Date(value);
			return Number.isNaN(parsed.getTime()) ? '—' : parsed.toLocaleString();
		},
		formatDate: (value: string | null | undefined) => {
			if (!value) return '—';
			const parsed = new Date(value);
			return Number.isNaN(parsed.getTime()) ? '—' : parsed.toLocaleDateString();
		},
		formatReviewQueueDate: (value: string | null | undefined) => {
			if (!value) return '—';
			const parsed = new Date(value);
			if (Number.isNaN(parsed.getTime())) return '—';
			const day = parsed.getUTCDate();
			const month = shortMonthNames[parsed.getUTCMonth()] ?? '—';
			const year = `${parsed.getUTCFullYear()}`.slice(-2);
			return `${day} ${month}, ${year}`;
		},
		async loadReviewQueue() {
			this.reviewQueueLoading = true;
			this.reviewQueueError = '';
			try {
				const response = await fetch('/v1/api/journals?reviewed=false&sort-by=created_at&sort-order=asc&limit=10');
				if (!response.ok) throw new Error('Failed to load review queue');
				const envelope = (await response.json()) as Envelope<JournalList>;
				this.reviewQueue = envelope.data?.journals ?? [];
			} catch (err) {
				this.reviewQueueError = err instanceof Error ? err.message : 'Unable to load review queue.';
			} finally {
				this.reviewQueueLoading = false;
			}
		},
	};
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalDetailPage', journalDetailPage);
});

export {};

import {
	type Journal,
	type Envelope,
	type JournalNote,
	type JournalNoteCreate,
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
		REJECTED: 'border-violet-300 bg-violet-50 text-violet-800',
	},
	type: {
		REJECTED: 'border-rose-300 bg-rose-50 text-rose-800',
		RESULT: 'border-emerald-300 bg-emerald-50 text-emerald-800',
		SET: 'border-indigo-300 bg-indigo-50 text-indigo-800',
	},
};

const normalizeTag = (value: string): string => (value ?? '').trim().toUpperCase();

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
		noteContent: '',
		reviewMessage: '',
		reviewMessageType: 'error',
		noteMessage: '',
		noteMessageType: 'error',
		normalizeStatus: normalizeTag,
		statusBadgeClass: (value: string) => badgeClassMap.status[normalizeTag(value)] ?? defaultBadgeClass,
		typeBadgeClass: (value: string) => badgeClassMap.type[normalizeTag(value)] ?? defaultBadgeClass,
		feedbackClass: (type: string) =>
			type === 'success' ? 'text-emerald-700' : 'text-rose-700',
		reviewToggleLabel(this: any) {
			return this.journal?.reviewed_at ? 'Mark Pending' : 'Mark Reviewed';
		},
		reviewButtonClass(this: any) {
			return this.journal?.reviewed_at
				? 'border-amber-300 bg-amber-50 text-amber-800 hover:bg-amber-100 focus:border-amber-400 focus:ring-amber-200'
				: 'border-emerald-300 bg-emerald-50 text-emerald-800 hover:bg-emerald-100 focus:border-emerald-400 focus:ring-emerald-200';
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
				const response = await fetch(`/v1/api/journals/${this.journalId}`, {
					method: 'PATCH',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(payload),
				});
				if (!response.ok) throw new Error(response.status === 404 ? 'Journal not found' : 'Failed to update review date');
				const envelope = (await response.json()) as Envelope<JournalReviewStatusResponse>;
				if (this.journal) {
					this.journal.reviewed_at = envelope.data.reviewed_at;
				}
				this.reviewMessageType = 'success';
				this.reviewMessage = envelope.data.reviewed_at ? 'Journal marked reviewed.' : 'Journal marked not reviewed.';
			} catch (err) {
				this.reviewMessage = err instanceof Error ? err.message : 'Unable to update review date.';
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
	};
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalDetailPage', journalDetailPage);
});

export {};

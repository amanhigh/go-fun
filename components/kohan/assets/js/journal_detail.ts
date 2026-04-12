import { type Journal, type Envelope } from './journal_models';
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
		normalizeStatus: normalizeTag,
		statusBadgeClass: (value: string) => badgeClassMap.status[normalizeTag(value)] ?? defaultBadgeClass,
		typeBadgeClass: (value: string) => badgeClassMap.type[normalizeTag(value)] ?? defaultBadgeClass,
		// Image helpers
		timeframeChipClass: image.chipClass,
		sortedImages: () => image.sorted(this.journal?.images),
		resolveImageSrc: image.resolve,
		previewImageSrc: () => image.resolve(this.previewImage()?.file_name ?? '', this.previewImage()?.created_at),
		previewImageLabel: () => image.label(this.previewImage()),
		previewImageCounter: () => image.counter(this.selectedImageIndex, this.sortedImages().length),
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
		hasError: () => this.errorMessage !== '',
		hasImagePreview: () => this.selectedImageIndex >= 0 && this.selectedImageIndex < this.sortedImages().length,
		openImagePreview: (index: number) => { this.selectedImageIndex = index; },
		closeImagePreview: () => { this.selectedImageIndex = -1; },
		canPrevImage: () => this.selectedImageIndex > 0,
		canNextImage: () => this.selectedImageIndex >= 0 && this.selectedImageIndex < this.sortedImages().length - 1,
		prevImage: (wrap = false) => {
			if (this.canPrevImage()) this.selectedImageIndex--;
			else if (wrap && this.sortedImages().length > 0) this.selectedImageIndex = this.sortedImages().length - 1;
		},
		nextImage: (wrap = false) => {
			if (this.canNextImage()) this.selectedImageIndex++;
			else if (wrap && this.sortedImages().length > 0) this.selectedImageIndex = 0;
		},
		previewImage: () => this.sortedImages()[this.selectedImageIndex] ?? null,
		previewImageTimeframe: () => this.previewImage()?.timeframe ?? '',
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

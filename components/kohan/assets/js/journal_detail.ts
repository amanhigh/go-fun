import { type Journal, type Envelope } from './journal_models';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalDetailPage>): void;
};

const defaultBadgeClass = 'border-slate-300 bg-slate-50 text-slate-700';

const statusBadgeClassMap: Record<string, string> = {
	SUCCESS: 'border-emerald-300 bg-emerald-50 text-emerald-800',
	FAIL: 'border-rose-300 bg-rose-50 text-rose-800',
	RUNNING: 'border-sky-300 bg-sky-50 text-sky-800',
	SET: 'border-amber-300 bg-amber-50 text-amber-800',
	REJECTED: 'border-violet-300 bg-violet-50 text-violet-800',
};

const typeBadgeClassMap: Record<string, string> = {
	REJECTED: 'border-rose-300 bg-rose-50 text-rose-800',
	RESULT: 'border-emerald-300 bg-emerald-50 text-emerald-800',
	SET: 'border-indigo-300 bg-indigo-50 text-indigo-800',
};

const normalizeTag = (value: string): string => (value ?? '').trim().toUpperCase();

function journalDetailPage() {
	return {
		journalId: '',
		journal: null as Journal | null,
		selectedImageIndex: -1,
		loading: false,
		errorMessage: '',
		normalizeStatus(value: string) {
			return normalizeTag(value);
		},
		statusBadgeClass(value: string) {
			return statusBadgeClassMap[normalizeTag(value)] ?? defaultBadgeClass;
		},
		timeframeChipClass(value: string) {
			const classes: Record<string, string> = {
				YR: 'border-fuchsia-400 bg-fuchsia-200 text-fuchsia-950',
				SMN: 'border-indigo-400 bg-indigo-200 text-indigo-950',
				TMN: 'border-cyan-400 bg-cyan-200 text-cyan-950',
				MN: 'border-emerald-400 bg-emerald-200 text-emerald-950',
				WK: 'border-amber-400 bg-amber-200 text-amber-950',
				DL: 'border-slate-400 bg-slate-200 text-slate-950',
			};
			return classes[normalizeTag(value)] ?? 'border-zinc-300 bg-zinc-100 text-zinc-900';
		},
		typeBadgeClass(value: string) {
			return typeBadgeClassMap[normalizeTag(value)] ?? defaultBadgeClass;
		},
		async loadJournal() {
			this.loading = true;
			this.errorMessage = '';
			try {
				const response = await fetch(`/v1/api/journals/${this.journalId}`);
				if (!response.ok) {
					if (response.status === 404) {
						throw new Error('Journal not found');
					}
					throw new Error('Failed to load journal');
				}
				const envelope = (await response.json()) as Envelope<Journal>;
				this.journal = envelope.data ?? null;
			} catch (err) {
				this.errorMessage = err instanceof Error ? err.message : 'Unable to load journal details. Please try again.';
			} finally {
				this.loading = false;
			}
		},
		hasError() {
			return this.errorMessage !== '';
		},
		timeframeRank(value: string) {
			const ranks: Record<string, number> = {
				YR: 600,
				SMN: 500,
				TMN: 400,
				MN: 300,
				WK: 200,
				DL: 100,
			};
			return ranks[normalizeTag(value)] ?? 0;
		},
		sortedImages() {
			const images = this.journal?.images ?? [];
			return [...images].sort((a, b) => this.timeframeRank(b.timeframe) - this.timeframeRank(a.timeframe));
		},
		hasImagePreview() {
			return this.selectedImageIndex >= 0 && this.selectedImageIndex < this.sortedImages().length;
		},
		openImagePreview(index: number) {
			this.selectedImageIndex = index;
		},
		closeImagePreview() {
			this.selectedImageIndex = -1;
		},
		canPrevImage() {
			return this.selectedImageIndex > 0;
		},
		canNextImage() {
			return this.selectedImageIndex >= 0 && this.selectedImageIndex < this.sortedImages().length - 1;
		},
		prevImage(wrap = false) {
			if (this.canPrevImage()) {
				this.selectedImageIndex -= 1;
				return;
			}
			if (wrap && this.sortedImages().length > 0) {
				this.selectedImageIndex = this.sortedImages().length - 1;
			}
		},
		nextImage(wrap = false) {
			if (this.canNextImage()) {
				this.selectedImageIndex += 1;
				return;
			}
			if (wrap && this.sortedImages().length > 0) {
				this.selectedImageIndex = 0;
			}
		},
		previewImage() {
			if (!this.hasImagePreview()) return null;
			return this.sortedImages()[this.selectedImageIndex] ?? null;
		},
		previewImageSrc() {
			const image = this.previewImage();
			if (!image) return '';
			return this.resolveImageSrc(image.file_name, image.created_at);
		},
		previewImageLabel() {
			const image = this.previewImage();
			if (!image) return '';
			return image.timeframe ? `${image.timeframe} • ${image.file_name}` : image.file_name;
		},
		previewImageTimeframe() {
			const image = this.previewImage();
			if (!image) return '';
			return image.timeframe ?? '';
		},
		resolveImageSrc(fileName: string, createdAt?: string) {
			if (!fileName) return '';
			if (fileName.startsWith('http://') || fileName.startsWith('https://') || fileName.startsWith('/')) {
				return fileName;
			}
			if (createdAt) {
				const date = new Date(createdAt);
				if (!Number.isNaN(date.getTime())) {
					const year = date.getFullYear();
					const month = String(date.getMonth() + 1).padStart(2, '0');
					return `/journal-images/${year}/${month}/${fileName}`;
				}
			}
			return '/journal-images/' + fileName;
		},
		formatTimestamp(value: string | null | undefined) {
			if (!value) return '—';
			const parsed = new Date(value);
			if (Number.isNaN(parsed.getTime())) return '—';
			return parsed.toLocaleString();
		},
	};
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalDetailPage', journalDetailPage);
});

export {};

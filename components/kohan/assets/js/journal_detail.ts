import { type Journal, type Envelope } from './journal_models';

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
	timeframe: {
		YR: 'border-fuchsia-400 bg-fuchsia-200 text-fuchsia-950',
		SMN: 'border-indigo-400 bg-indigo-200 text-indigo-950',
		TMN: 'border-cyan-400 bg-cyan-200 text-cyan-950',
		MN: 'border-emerald-400 bg-emerald-200 text-emerald-950',
		WK: 'border-amber-400 bg-amber-200 text-amber-950',
		DL: 'border-slate-400 bg-slate-200 text-slate-950',
	},
};

const timeframeRankMap: Record<string, number> = { YR: 600, SMN: 500, TMN: 400, MN: 300, WK: 200, DL: 100 };

const normalizeTag = (value: string): string => (value ?? '').trim().toUpperCase();

function journalDetailPage() {
	return {
		journalId: '',
		journal: null as Journal | null,
		selectedImageIndex: -1,
		loading: false,
		errorMessage: '',
		normalizeStatus: normalizeTag,
		statusBadgeClass: (value: string) => badgeClassMap.status[normalizeTag(value)] ?? defaultBadgeClass,
		timeframeChipClass: (value: string) => badgeClassMap.timeframe[normalizeTag(value)] ?? 'border-zinc-300 bg-zinc-100 text-zinc-900',
		typeBadgeClass: (value: string) => badgeClassMap.type[normalizeTag(value)] ?? defaultBadgeClass,
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
		timeframeRank: (value: string) => timeframeRankMap[normalizeTag(value)] ?? 0,
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
				this.selectedImageIndex--;
			} else if (wrap && this.sortedImages().length > 0) {
				this.selectedImageIndex = this.sortedImages().length - 1;
			}
		},
		nextImage(wrap = false) {
			if (this.canNextImage()) {
				this.selectedImageIndex++;
			} else if (wrap && this.sortedImages().length > 0) {
				this.selectedImageIndex = 0;
			}
		},
		previewImage() {
			return this.sortedImages()[this.selectedImageIndex] ?? null;
		},
		previewImageSrc() {
			const image = this.previewImage();
			return image ? this.resolveImageSrc(image.file_name, image.created_at) : '';
		},
		previewImageLabel() {
			const image = this.previewImage();
			return image ? (image.timeframe ? `${image.timeframe} • ${image.file_name}` : image.file_name) : '';
		},
		previewImageTimeframe() {
			return this.previewImage()?.timeframe ?? '';
		},
		resolveImageSrc(fileName: string, createdAt?: string) {
			if (!fileName || fileName.startsWith('http://') || fileName.startsWith('https://') || fileName.startsWith('/')) {
				return fileName || '';
			}
			if (!createdAt) return '/journal-images/' + fileName;
			const date = new Date(createdAt);
			if (Number.isNaN(date.getTime())) return '/journal-images/' + fileName;
			return `/journal-images/${date.getFullYear()}/${String(date.getMonth() + 1).padStart(2, '0')}/${fileName}`;
		},
		formatTimestamp(value: string | null | undefined) {
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

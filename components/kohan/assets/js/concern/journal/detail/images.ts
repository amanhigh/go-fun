import type { JournalImage } from '../../../types/journal_api';
import type { JournalDetailPageProvider, JournalImagesConcern } from '../../../types/journal_detail_concern';

export interface ImageHelper {
	chipClass(timeframe: string): string;
	sorted(images: JournalImage[] | undefined): JournalImage[];
	resolve(fileName: string, createdAt?: string): string;
	label(image: JournalImage | null | undefined): string;
	counter(current: number, total: number): string;
}

const chipClassMap: Record<string, string> = {
	YR: 'border-fuchsia-400 bg-fuchsia-200 text-fuchsia-950',
	SMN: 'border-sky-400 bg-sky-200 text-sky-950',
	TMN: 'border-emerald-400 bg-emerald-200 text-emerald-950',
	MN: 'border-orange-400 bg-orange-200 text-orange-950',
	WK: 'border-yellow-400 bg-yellow-200 text-yellow-950',
	DL: 'border-slate-400 bg-slate-200 text-slate-950',
};

const rankMap: Record<string, number> = { YR: 600, SMN: 500, TMN: 400, MN: 300, WK: 200, DL: 100 };

const normalize = (value: string): string => (value ?? '').trim().toUpperCase();

export function createImageHelper(): ImageHelper {
	return {
		chipClass(timeframe) {
			return chipClassMap[normalize(timeframe)] ?? 'border-zinc-300 bg-zinc-100 text-zinc-900';
		},
		sorted(images) {
			if (!images?.length) return [];
			return [...images].sort((a, b) => {
				const aDate = a.created_at ? new Date(a.created_at).getTime() : Number.POSITIVE_INFINITY;
				const bDate = b.created_at ? new Date(b.created_at).getTime() : Number.POSITIVE_INFINITY;
				if (aDate !== bDate) return aDate - bDate;
				const timeframeDiff = (rankMap[normalize(b.timeframe)] ?? 0) - (rankMap[normalize(a.timeframe)] ?? 0);
				if (timeframeDiff !== 0) return timeframeDiff;
				return normalize(a.file_name).localeCompare(normalize(b.file_name));
			});
		},
		resolve(fileName, createdAt) {
			if (!fileName) return '';
			if (fileName.startsWith('http://') || fileName.startsWith('https://') || fileName.startsWith('/')) return fileName;
			if (!createdAt) return '/journal/images/' + fileName;
			const date = new Date(createdAt);
			if (Number.isNaN(date.getTime())) return '/journal/images/' + fileName;
			return `/journal/images/${date.getFullYear()}/${String(date.getMonth() + 1).padStart(2, '0')}/${fileName}`;
		},
		label(image) {
			if (!image) return '';
			return image.timeframe ? `${image.timeframe} • ${image.file_name}` : image.file_name;
		},
		counter(current, total) {
			return `${current + 1} / ${total}`;
		},
	};
}

export function NewImagesConcern(pg: JournalDetailPageProvider, image: ImageHelper): JournalImagesConcern {
	return {
		resolveImageSrc: image.resolve,
		timeframeChipClass: image.chipClass,
		countLabel(this: any) {
			const count = this.sorted().length;
			return `${count} timeframe image${count === 1 ? '' : 's'}`;
		},
		sorted() {
			return image.sorted(pg().current.journal?.images);
		},
		tileTitle(this: any, imageItem: JournalImage) {
			return imageItem.file_name;
		},
		tileSrc(this: any, imageItem: JournalImage) {
			return image.resolve(imageItem.file_name, imageItem.created_at);
		},
		tileAlt(this: any, imageItem: JournalImage) {
			return image.label(imageItem);
		},
	};
}

import type { JournalImage } from './journal_models';

export interface ImageHelper {
	chipClass(timeframe: string): string;
	sorted(images: JournalImage[] | undefined): JournalImage[];
	resolve(fileName: string, createdAt?: string): string;
	label(image: JournalImage | null | undefined): string;
	counter(current: number, total: number): string;
}

const chipClassMap: Record<string, string> = {
	YR: 'border-fuchsia-400 bg-fuchsia-200 text-fuchsia-950',
	SMN: 'border-indigo-400 bg-indigo-200 text-indigo-950',
	TMN: 'border-cyan-400 bg-cyan-200 text-cyan-950',
	MN: 'border-emerald-400 bg-emerald-200 text-emerald-950',
	WK: 'border-amber-400 bg-amber-200 text-amber-950',
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
			return [...images].sort((a, b) => (rankMap[normalize(b.timeframe)] ?? 0) - (rankMap[normalize(a.timeframe)] ?? 0));
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

import type { JournalImage } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

const TIMEFRAME_RANK: Record<string, number> = { YR: 600, SMN: 500, TMN: 400, MN: 300, WK: 200, DL: 100 };

const normalize = (value: string): string => (value ?? '').trim().toUpperCase();

export function NewImagesConcern(pg: JournalDetailPageProvider) {
	return {
		sorted(): JournalImage[] {
			const images = pg().current.journal?.images;
			if (!images?.length) return [];
			return [...images].sort((a, b) => {
				const aDate = a.created_at ? new Date(a.created_at).getTime() : Number.POSITIVE_INFINITY;
				const bDate = b.created_at ? new Date(b.created_at).getTime() : Number.POSITIVE_INFINITY;
				if (aDate !== bDate) return aDate - bDate;
				const timeframeDiff = (TIMEFRAME_RANK[normalize(b.timeframe)] ?? 0) - (TIMEFRAME_RANK[normalize(a.timeframe)] ?? 0);
				if (timeframeDiff !== 0) return timeframeDiff;
				return normalize(a.file_name).localeCompare(normalize(b.file_name));
			});
		},

		src(image: JournalImage): string {
			if (!image.file_name) return '';
			if (image.file_name.startsWith('http://') || image.file_name.startsWith('https://') || image.file_name.startsWith('/')) return image.file_name;
			if (!image.created_at) return '/journal/images/' + image.file_name;
			const date = new Date(image.created_at);
			if (Number.isNaN(date.getTime())) return '/journal/images/' + image.file_name;
			return `/journal/images/${date.getFullYear()}/${String(date.getMonth() + 1).padStart(2, '0')}/${image.file_name}`;
		},

		label(image: JournalImage | null | undefined): string {
			if (!image) return '';
			return image.timeframe ? `${image.timeframe} • ${image.file_name}` : image.file_name;
		},

		countLabel(this: any): string {
			const count = this.sorted().length;
			return `${count} timeframe image${count === 1 ? '' : 's'}`;
		},
	};
}

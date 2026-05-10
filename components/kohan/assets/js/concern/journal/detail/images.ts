import type { JournalImage } from '../../../types/api/journal/response';
import type { JournalTimeframe } from '../../../types/api/journal/enums';
import type { JournalImageView, JournalDetailPageProvider } from '../../../types/journal/detail';

const TIMEFRAME_RANK: Record<JournalTimeframe, number> = { YR: 600, SMN: 500, TMN: 400, MN: 300, WK: 200, DL: 100 };

function toImageView(image: JournalImage): JournalImageView {
	return {
		...image,
		src: imageSrc(image),
		label: imageLabel(image),
	};
}

function imageSrc(image: JournalImage): string {
	if (!image.file_name) return '';
	if (image.file_name.startsWith('http://') || image.file_name.startsWith('https://') || image.file_name.startsWith('/')) return image.file_name;
	const path = (() => {
		if (!image.created_at) return '/journal/images/' + image.file_name;
		const date = new Date(image.created_at);
		if (Number.isNaN(date.getTime())) return '/journal/images/' + image.file_name;
		return `/journal/images/${date.getFullYear()}/${String(date.getMonth() + 1).padStart(2, '0')}/${image.file_name}`;
	})();
	// DONOT REMOVE: Cache-busting query param — prevents Brave/Firefox from serving stale partial screenshots
	// IT WORKED. DO NOT DELETE. Append created_at so browser re-fetches after screenshot updates.
	return path + '?t=' + encodeURIComponent(image.created_at || image.file_name);
}

function imageLabel(image: JournalImage): string {
	return image.timeframe ? `${image.timeframe} • ${image.file_name}` : image.file_name;
}

function compareImages(a: JournalImage, b: JournalImage): number {
	const aDate = a.created_at ? new Date(a.created_at).getTime() : Number.POSITIVE_INFINITY;
	const bDate = b.created_at ? new Date(b.created_at).getTime() : Number.POSITIVE_INFINITY;
	if (aDate !== bDate) return aDate - bDate;
	const timeframeDiff = (TIMEFRAME_RANK[b.timeframe] ?? 0) - (TIMEFRAME_RANK[a.timeframe] ?? 0);
	if (timeframeDiff !== 0) return timeframeDiff;
	return a.file_name.localeCompare(b.file_name);
}

export function NewImagesConcern(pg: JournalDetailPageProvider) {
	return {
		sorted(): JournalImageView[] {
			const images = pg().journal.detail!.images;
			if (!images?.length) return [];
			return [...images].map(toImageView).sort(compareImages);
		},

		countLabel(): string {
			const count = pg().journal.detail!.images.length;
			return `${count} timeframe image${count === 1 ? '' : 's'}`;
		},
	};
}

import type { JournalImage, JournalTimeframe } from '../../../types/journal_api';
import type { JournalImageView } from '../../../types/journal_detail_concern';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

const TIMEFRAME_RANK: Record<JournalTimeframe, number> = { YR: 600, SMN: 500, TMN: 400, MN: 300, WK: 200, DL: 100 };

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
			const images = pg().current.journal?.images;
			if (!images?.length) return [];
			return [...images].sort(compareImages) as JournalImageView[];
		},

		countLabel(): string {
			const count = pg().current.journal?.images?.length ?? 0;
			return `${count} timeframe image${count === 1 ? '' : 's'}`;
		},
	};
}

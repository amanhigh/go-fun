import type { JournalImage } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewImagePreviewConcern(pg: JournalDetailPageProvider) {
	return {
		index: -1,

		open(idx: number) { pg().preview.index = idx; },
		close() { pg().preview.index = -1; },

		current(): JournalImage | null {
			return pg().images.sorted()[pg().preview.index] ?? null;
		},
		timeframe() { return this.current()?.timeframe ?? ''; },
		src() {
			const img = this.current();
			return img ? pg().images.resolveImageSrc(img.file_name, img.created_at) : '';
		},
		label() {
			const img = this.current();
			return img ? `${img.timeframe} • ${img.file_name}` : '';
		},
		counter() {
			const total = pg().images.sorted().length;
			return `${pg().preview.index + 1} / ${total}`;
		},

		hasPreview() {
			const idx = pg().preview.index;
			return idx >= 0 && idx < pg().images.sorted().length;
		},
		canPrev() { return pg().preview.index > 0; },
		canNext() {
			const idx = pg().preview.index;
			return idx >= 0 && idx < pg().images.sorted().length - 1;
		},
		prev(this: any, wrap = false) {
			const total = pg().images.sorted().length;
			if (this.canPrev()) pg().preview.index--;
			else if (wrap && total > 0) pg().preview.index = total - 1;
		},
		next(this: any, wrap = false) {
			const total = pg().images.sorted().length;
			if (this.canNext()) pg().preview.index++;
			else if (wrap && total > 0) pg().preview.index = 0;
		},
	};
}

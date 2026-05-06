import type { JournalImage } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewImagePreviewConcern(pg: JournalDetailPageProvider) {
	return {
		index: -1,

		open(idx: number) { pg().preview.index = idx; },
		close() { pg().preview.index = -1; },

		current(this: any): JournalImage | null {
			return this.images?.sorted()?.[pg().preview.index] ?? null;
		},
		timeframe(this: any) { return this.current()?.timeframe ?? ''; },
		src(this: any) {
			const img = this.current();
			return img ? pg().images.resolveImageSrc(img.file_name, img.created_at) : '';
		},
		label(this: any) {
			const img = this.current();
			return img ? `${img.timeframe} • ${img.file_name}` : '';
		},
		counter(this: any) {
			return `${pg().preview.index + 1} / ${this.images?.sorted()?.length ?? 0}`;
		},

		hasPreview(this: any) {
			const idx = pg().preview.index;
			return idx >= 0 && idx < (this.images?.sorted()?.length ?? 0);
		},
		canPrev() { return pg().preview.index > 0; },
		canNext(this: any) {
			const idx = pg().preview.index;
			return idx >= 0 && idx < (this.images?.sorted()?.length ?? 0) - 1;
		},
		prev(this: any, wrap = false) {
			const total = this.images?.sorted()?.length ?? 0;
			if (this.canPrev()) pg().preview.index--;
			else if (wrap && total > 0) pg().preview.index = total - 1;
		},
		next(this: any, wrap = false) {
			const total = this.images?.sorted()?.length ?? 0;
			if (this.canNext()) pg().preview.index++;
			else if (wrap && total > 0) pg().preview.index = 0;
		},
	};
}

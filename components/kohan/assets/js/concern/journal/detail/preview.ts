import type { JournalImageView, JournalDetailPageProvider } from '../../../types/journal/detail';

function currentImage(pg: JournalDetailPageProvider, index: number): JournalImageView | null {
	const images = pg().images.sorted();
	return images[index] ?? null;
}

export function NewPreviewConcern(pg: JournalDetailPageProvider) {
	return {
		index: -1,

		open(idx: number) { this.index = idx; },
		close() { this.index = -1; },

		timeframe() { return currentImage(pg, this.index)?.timeframe ?? ''; },
		src() { return currentImage(pg, this.index)?.src ?? ''; },
		label() { return currentImage(pg, this.index)?.label ?? ''; },
		fileName() { return currentImage(pg, this.index)?.file_name ?? ''; },

		counter() {
			const img = currentImage(pg, this.index);
			if (!img) return '';
			const total = pg().images.sorted().length;
			return `${this.index + 1} / ${total}`;
		},

		hasPreview() {
			const images = pg().images.sorted();
			return this.index >= 0 && this.index < images.length;
		},
		prev() {
			if (this.index > 0) this.index--;
		},
		next() {
			const images = pg().images.sorted();
			if (this.index >= 0 && this.index < images.length - 1) this.index++;
		},
		wrapPrev() {
			const total = pg().images.sorted().length;
			if (total > 0) this.index = this.index > 0 ? this.index - 1 : total - 1;
		},
		wrapNext() {
			const total = pg().images.sorted().length;
			if (total > 0) this.index = this.index < total - 1 ? this.index + 1 : 0;
		},
	};
}

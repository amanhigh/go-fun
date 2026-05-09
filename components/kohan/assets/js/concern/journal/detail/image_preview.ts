import type { JournalImageView } from '../../../types/journal_detail_concern';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewImagePreviewConcern(pg: JournalDetailPageProvider) {
	return {
		index: -1,

		open(idx: number) { this.index = idx; },
		close() { this.index = -1; },

		current(): JournalImageView | null {
			const images = pg().images.sorted();
			return images[this.index] ?? null;
		},
		timeframe() { return this.current()?.timeframe ?? ''; },
		src() { return this.current()?.src ?? ''; },
		label() { return this.current()?.label ?? ''; },
		counter() {
			const total = pg().images.sorted().length;
			return `${this.index + 1} / ${total}`;
		},

		hasPreview() {
			const images = pg().images.sorted();
			return this.index >= 0 && this.index < images.length;
		},
		canPrev() { return this.index > 0; },
		canNext() {
			const images = pg().images.sorted();
			return this.index >= 0 && this.index < images.length - 1;
		},
		prev() {
			if (this.canPrev()) this.index--;
		},
		next() {
			if (this.canNext()) this.index++;
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

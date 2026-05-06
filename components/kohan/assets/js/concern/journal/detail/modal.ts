import type { JournalImage } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';
import type { ImageHelper } from './images';

export function NewPreviewConcern(pg: JournalDetailPageProvider, image: ImageHelper) {
	return {
		timeframeChipClass: image.chipClass,
		imageCountLabel() {
			const count = this.sortedImages().length;
			return `${count} timeframe image${count === 1 ? '' : 's'}`;
		},
		sortedImages() {
			return image.sorted(pg().journal?.images);
		},
		imageTileTitle(this: any, imageItem: JournalImage) {
			return imageItem.file_name;
		},
		imageTileSrc(this: any, imageItem: JournalImage) {
			return image.resolve(imageItem.file_name, imageItem.created_at);
		},
		imageTileAlt(this: any, imageItem: JournalImage) {
			return image.label(imageItem);
		},
		resolveImageSrc: image.resolve,
		previewImageSrc() {
			return image.resolve(this.previewImage()?.file_name ?? '', this.previewImage()?.created_at);
		},
		previewImageLabel() {
			return image.label(this.previewImage());
		},
		previewImageCounter() {
			return image.counter(pg().selectedImageIndex, this.sortedImages().length);
		},
		hasImagePreview() {
			return pg().selectedImageIndex >= 0 && pg().selectedImageIndex < this.sortedImages().length;
		},
		openImagePreview(index: number) {
			pg().selectedImageIndex = index;
		},
		closeImagePreview() {
			pg().selectedImageIndex = -1;
		},
		canPrevImage() {
			return pg().selectedImageIndex > 0;
		},
		canNextImage() {
			return pg().selectedImageIndex >= 0 && pg().selectedImageIndex < this.sortedImages().length - 1;
		},
		prevImage(wrap = false) {
			if (this.canPrevImage()) pg().selectedImageIndex--;
			else if (wrap && this.sortedImages().length > 0) pg().selectedImageIndex = this.sortedImages().length - 1;
		},
		nextImage(wrap = false) {
			if (this.canNextImage()) pg().selectedImageIndex++;
			else if (wrap && this.sortedImages().length > 0) pg().selectedImageIndex = 0;
		},
		previewImage(): JournalImage | null {
			return this.sortedImages()[pg().selectedImageIndex] ?? null;
		},
		previewImageTimeframe() {
			return this.previewImage()?.timeframe ?? '';
		},
	};
}

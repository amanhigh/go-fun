import type { JournalImage } from './journal_models';
import type { ImageHelper } from './journal_images';

export function createJournalDetailPreview(image: ImageHelper) {
	return {
		timeframeChipClass: image.chipClass,
		sortedImages(this: any) {
			return image.sorted(this.journal?.images);
		},
		resolveImageSrc: image.resolve,
		previewImageSrc(this: any) {
			return image.resolve(this.previewImage()?.file_name ?? '', this.previewImage()?.created_at);
		},
		previewImageLabel(this: any) {
			return image.label(this.previewImage());
		},
		previewImageCounter(this: any) {
			return image.counter(this.selectedImageIndex, this.sortedImages().length);
		},
		hasImagePreview(this: any) {
			return this.selectedImageIndex >= 0 && this.selectedImageIndex < this.sortedImages().length;
		},
		openImagePreview(this: any, index: number) {
			this.selectedImageIndex = index;
		},
		closeImagePreview(this: any) {
			this.selectedImageIndex = -1;
		},
		canPrevImage(this: any) {
			return this.selectedImageIndex > 0;
		},
		canNextImage(this: any) {
			return this.selectedImageIndex >= 0 && this.selectedImageIndex < this.sortedImages().length - 1;
		},
		prevImage(this: any, wrap = false) {
			if (this.canPrevImage()) this.selectedImageIndex--;
			else if (wrap && this.sortedImages().length > 0) this.selectedImageIndex = this.sortedImages().length - 1;
		},
		nextImage(this: any, wrap = false) {
			if (this.canNextImage()) this.selectedImageIndex++;
			else if (wrap && this.sortedImages().length > 0) this.selectedImageIndex = 0;
		},
		previewImage(this: any): JournalImage | null {
			return this.sortedImages()[this.selectedImageIndex] ?? null;
		},
		previewImageTimeframe(this: any) {
			return this.previewImage()?.timeframe ?? '';
		},
	};
}

import { getErrorMessage } from '../../../lib/error';
import type { JournalImage } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function imageSrc(image: JournalImage): string {
	if (!image.file_name) return '';
	if (image.file_name.startsWith('http://') || image.file_name.startsWith('https://') || image.file_name.startsWith('/')) return image.file_name;
	if (!image.created_at) return '/journal/images/' + image.file_name;
	const date = new Date(image.created_at);
	if (Number.isNaN(date.getTime())) return '/journal/images/' + image.file_name;
	return `/journal/images/${date.getFullYear()}/${String(date.getMonth() + 1).padStart(2, '0')}/${image.file_name}`;
}

function imageLabel(image: JournalImage): string {
	return image.timeframe ? `${image.timeframe} • ${image.file_name}` : image.file_name;
}

function normalizeJournal(journal: any) {
	if (!journal) return null;
	return {
		...journal,
		images: (journal.images ?? []).map((img: JournalImage) => ({
			...img,
			src: imageSrc(img),
			label: imageLabel(img),
		})),
		tags: journal.tags ?? [],
		notes: journal.notes ?? [],
	};
}

export function NewCurrentJournalConcern(pg: JournalDetailPageProvider) {
	return {
		journalId: '',
		journal: null,
		loading: true,
		errorMessage: '',

		hasError(this: any) { return this.errorMessage !== ''; },

		async loadJournal(this: any) {
			this.loading = true;
			this.errorMessage = '';
			try {
				const envelope = await pg().client.get(this.journalId);
				this.journal = normalizeJournal(envelope.data);
				pg().sidebar.tags.sync(this.journal?.tags);
				pg().sidebar.notes.sync(this.journal?.notes);
			} catch (err) {
				this.errorMessage = getErrorMessage(err, 'Unable to load journal details. Please try again.');
			} finally {
				this.loading = false;
			}
		},
	};
}

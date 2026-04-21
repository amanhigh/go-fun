import { NewJournalClient, type Journal } from '../client/journal';
import { NewJournalNoteClient } from '../client/journal_note';
import { NewJournalTagClient } from '../client/journal_tag';
import { createJournalDetailFormatters } from './formatters';
import { createImageHelper } from '../journal_images';
import { createJournalDetailNotes } from './notes';
import { createJournalDetailPreview } from './preview';
import { createJournalDetailReview } from './review';
import { createJournalDetailTags, managementTagPresets } from './tags';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalDetailPage>): void;
};

function journalDetailPage() {
	const journalClient = NewJournalClient();
	const noteClient = NewJournalNoteClient();
	const tagClient = NewJournalTagClient();
	const image = createImageHelper();

	return {
		journalId: '',
		journal: null as Journal | null,
		selectedImageIndex: -1,
		loading: false,
		errorMessage: '',
		reviewSubmitting: false,
		noteSubmitting: false,
		noteDeletingId: '' as string,
		noteContent: '',
		reviewMessage: '',
		reviewMessageType: 'error',
		noteMessage: '',
		noteMessageType: 'error',
		reviewQueue: [] as Journal[],
		reviewQueueLoading: false,
		reviewQueueError: '',
		managementTagPresets,
		managementTagSubmitting: false,
		managementTagPendingValue: '',
		managementTagMessage: '',
		managementTagMessageType: 'error',
		reasonTagInput: '',
		reasonTagOverride: '',
		reasonTagSubmitting: false,
		tagDeletingId: '',
		reasonTagMessage: '',
		reasonTagMessageType: 'error',
		...createJournalDetailFormatters(),
		...createJournalDetailPreview(image),
		...createJournalDetailReview(journalClient),
		...createJournalDetailNotes(noteClient),
		...createJournalDetailTags(tagClient),
		async loadJournal(this: any) {
			this.loading = true;
			this.errorMessage = '';
			try {
				const envelope = await journalClient.get(this.journalId);
				this.journal = envelope.data ?? null;
			} catch (err) {
				this.errorMessage = err instanceof Error ? err.message : 'Unable to load journal details. Please try again.';
			} finally {
				this.loading = false;
			}
		},
		hasError(this: any) {
			return this.errorMessage !== '';
		},
	};
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalDetailPage', journalDetailPage);
});

export {};

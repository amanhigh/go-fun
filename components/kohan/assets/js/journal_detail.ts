import type { Envelope, Journal } from './journal_models';
import { createJournalDetailFormatters } from './journal_detail_formatters';
import { createImageHelper } from './journal_images';
import { createJournalDetailNotes } from './journal_detail_notes';
import { createJournalDetailPreview } from './journal_detail_preview';
import { createJournalDetailReview } from './journal_detail_review';
import { createJournalDetailTags, managementTagPresets } from './journal_detail_tags';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalDetailPage>): void;
};

function journalDetailPage() {
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
		...createJournalDetailReview(),
		...createJournalDetailNotes(),
		...createJournalDetailTags(),
		async loadJournal(this: any) {
			this.loading = true;
			this.errorMessage = '';
			try {
				const response = await fetch(`/v1/api/journals/${this.journalId}`);
				if (!response.ok) throw new Error(response.status === 404 ? 'Journal not found' : 'Failed to load journal');
				const envelope = (await response.json()) as Envelope<Journal>;
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

import { NewJournalClient } from '../client/journal';
import { NewJournalNoteClient } from '../client/journal_note';
import { NewJournalTagClient } from '../client/journal_tag';
import type { Journal } from '../types/journal';
import type { JournalClient } from '../client/journal';
import type { JournalDetailPageState } from '../types/journal_detail_state';
import { getErrorMessage } from '../shared/error';
import { createJournalDetailFormatters } from '../concern/journal/detail/header';
import { createImageHelper } from '../concern/journal/detail/images';
import { createJournalDetailPreview } from '../concern/journal/detail/modal';
import { createSideBar } from '../concern/journal/detail/sidebar';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalDetailPage>): void;
};

function createEmptyJournal(): Journal {
	return {
		id: '',
		ticker: '',
		sequence: '',
		type: '',
		status: '',
		created_at: '',
		reviewed_at: null,
		images: [],
		tags: [],
		notes: [],
		deleted_at: null,
	};
}

function createJournalDetailPageState(): JournalDetailPageState {
	return {
		journalId: '',
		journal: createEmptyJournal(),
		selectedImageIndex: -1,
		loading: true,
		journalDeleting: false,
		errorMessage: '',
	};
}

function normalizeJournal(journal: any) {
	if (!journal) return null;

	return {
		...journal,
		images: journal.images ?? [],
		tags: journal.tags ?? [],
		notes: [...(journal.notes ?? [])].sort((left, right) => {
			const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
			const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
			return rightTime - leftTime;
		}),
	};
}

function createJournalDetailPageActions(journalClient: JournalClient) {
	return {
		syncSideBarCollections(this: any) {
			this.sidebar?.syncNotes?.(this.journal?.notes);
			this.sidebar?.syncTags?.(this.journal?.tags);
		},
		async loadJournal(this: any) {
			this.loading = true;
			this.errorMessage = '';
			try {
				const envelope = await journalClient.get(this.journalId);
				this.journal = normalizeJournal(envelope.data);
				this.syncSideBarCollections();
			} catch (err) {
				this.errorMessage = getErrorMessage(err, 'Unable to load journal details. Please try again.');
			} finally {
				this.loading = false;
			}
		},
		async deleteJournal(this: any) {
			if (!this.journal || this.journalDeleting) return;
			if (typeof window !== 'undefined' && !window.confirm('Delete this journal? This cannot be undone.')) return;
			this.journalDeleting = true;
			try {
				await journalClient.delete(this.journalId);
				window.location.href = '/journal';
			} catch (err) {
				this.errorMessage = getErrorMessage(err, 'Unable to delete journal.');
			} finally {
				this.journalDeleting = false;
			}
		},
		hasError(this: any) { return this.errorMessage !== ''; },
	};
}

function journalDetailPage() {
	const root = document.querySelector<HTMLElement>('[data-journal-detail-page]');
	const journalId = root?.dataset.journalId ?? '';
	const actionOpenStorageKey = root?.dataset.actionOpenStorageKey ?? '';
	const reviewModeStorageKey = root?.dataset.reviewModeStorageKey ?? '';
	const journalClient = NewJournalClient();
	const noteClient = NewJournalNoteClient();
	const tagClient = NewJournalTagClient();
	const image = createImageHelper();

	const state = createJournalDetailPageState();
	const formatters = createJournalDetailFormatters();
	const pageActions = createJournalDetailPageActions(journalClient);
	const preview = createJournalDetailPreview(image);

	const page: any = {
		...state,
		...formatters,
		...pageActions,
		...preview,
		init(this: any) {
			this.journalId = journalId;
			this.sidebar.initSidebarUiState(actionOpenStorageKey, reviewModeStorageKey);
			void this.loadJournal();
			void this.sidebar.loadReviewQueue();
		},
	};

	page.sidebar = createSideBar(page, journalClient, noteClient, tagClient);

	return page;
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalDetailPage', journalDetailPage);
});

export {};

import { NewJournalClient } from '../client/journal';
import { NewJournalNoteClient } from '../client/journal_note';
import { NewJournalTagClient } from '../client/journal_tag';
import type { JournalDetailPageData } from '../types/journal_detail_concern';
import { getErrorMessage } from '../shared/error';
import { NewPresentationConcern } from '../concern/journal/common/presentation';
import { NewHeaderConcern } from '../concern/journal/detail/header';
import { createImageHelper } from '../concern/journal/detail/images';
import { NewPreviewConcern } from '../concern/journal/detail/modal';
import { NewSidebarConcern } from '../concern/journal/detail/sidebar';
import '../types/platform';

function normalizeJournal(journal: any) {
	if (!journal) return null;

	return {
		...journal,
		images: journal.images ?? [],
		tags: journal.tags ?? [],
		notes: [...(journal.notes ?? [])].sort((left: any, right: any) => {
			const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
			const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
			return rightTime - leftTime;
		}),
	};
}

function journalDetailPage() {
	let page = {} as JournalDetailPageData;
	const pg = () => page;

	const journalClient = NewJournalClient();
	const noteClient = NewJournalNoteClient();
	const tagClient = NewJournalTagClient();
	const image = createImageHelper();

	// Clients
	page.client = journalClient;
	page.noteClient = noteClient;
	page.tagClient = tagClient;

	// State defaults
	page.journalId = '';
	page.journal = null;
	page.selectedImageIndex = -1;
	page.loading = true;
	page.journalDeleting = false;
	page.errorMessage = '';

	// Presentation
	page.presentation = NewPresentationConcern();

	// Header formatters (pure, no pg needed)
	Object.assign(page, NewHeaderConcern());

	// Preview/modal concern
	Object.assign(page, NewPreviewConcern(pg, image));

	// Sidebar nested concern
	page.sidebar = NewSidebarConcern(pg);

	// Page-level actions
	page.syncSideBarCollections = function syncSideBarCollections(this: any) {
		this.sidebar?.syncNotes?.(pg().journal?.notes);
		this.sidebar?.syncTags?.(pg().journal?.tags);
	};

	page.loadJournal = async function loadJournal(this: any) {
		this.loading = true;
		this.errorMessage = '';
		try {
			const envelope = await pg().client.get(pg().journalId);
			pg().journal = normalizeJournal(envelope.data);
			this.syncSideBarCollections();
		} catch (err) {
			this.errorMessage = getErrorMessage(err, 'Unable to load journal details. Please try again.');
		} finally {
			this.loading = false;
		}
	};

	page.deleteJournal = async function deleteJournal(this: any) {
		if (!pg().journal || this.journalDeleting) return;
		if (!window.confirm('Delete this journal? This cannot be undone.')) return;
		this.journalDeleting = true;
		try {
			await pg().client.delete(pg().journalId);
			window.location.href = '/journal';
		} catch (err) {
			this.errorMessage = getErrorMessage(err, 'Unable to delete journal.');
		} finally {
			this.journalDeleting = false;
		}
	};

	page.hasError = function hasError(this: any) { return this.errorMessage !== ''; };

	page.init = function init(this: any) {
		page = this;
		const el = this.$el as HTMLElement;
		this.journalId = el.dataset.journalId ?? '';
		this.sidebar.initSidebarUiState(
			el.dataset.actionOpenStorageKey ?? '',
			el.dataset.reviewModeStorageKey ?? '',
		);
		void this.loadJournal();
		void this.sidebar.loadReviewQueue();
	};

	return page;
}

document.addEventListener('alpine:init', () => {
	window.Alpine.data('journalDetailPage', journalDetailPage);
});

export {};

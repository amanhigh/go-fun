import { NewJournalClient } from '../client/journal';
import { NewJournalNoteClient } from '../client/journal_note';
import { NewJournalTagClient } from '../client/journal_tag';
import type { JournalDetailPage } from '../types/journal/detail';
import { NewPresentationConcern } from '../concern/present/factory';
import { NewJournalConcern } from '../concern/journal/detail/journal';
import { NewHeaderConcern } from '../concern/journal/detail/header';
import { NewImagesConcern } from '../concern/journal/detail/images';
import { NewPreviewConcern } from '../concern/journal/detail/preview';
import { NewSidebarConcern } from '../concern/journal/sidebar';
import '../types/core/platform';

function createJournalDetailPageData(journalId = '') {
	let page = {} as JournalDetailPage;
	const pg = () => page;

	page.client = NewJournalClient();
	page.noteClient = NewJournalNoteClient();
	page.tagClient = NewJournalTagClient();

	page.present = NewPresentationConcern();
	page.current = NewJournalConcern(pg);
	page.header = NewHeaderConcern(pg);
	page.images = NewImagesConcern(pg);
	page.preview = NewPreviewConcern(pg);
	page.sidebar = NewSidebarConcern(pg);

	page.init = function init(this: any) {
		page = this;

		this.current.journalId = journalId;
		this.sidebar.state.restorePersistedSidebarState();

		void this.current.loadJournal();
		void this.sidebar.reviewQueue.load();
	};

	return page;
}

document.addEventListener('alpine:init', () => {
	window.Alpine.data('journalDetailPage', createJournalDetailPageData);
});

export {};

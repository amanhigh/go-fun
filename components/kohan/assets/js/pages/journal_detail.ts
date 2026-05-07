import { NewJournalClient } from '../client/journal';
import { NewJournalNoteClient } from '../client/journal_note';
import { NewJournalTagClient } from '../client/journal_tag';
import type { JournalDetailPageData } from '../types/journal_detail_concern';
import { NewPresentationConcern } from '../concern/journal/common/presentation';
import { NewCurrentJournalConcern } from '../concern/journal/detail/current_journal';
import { NewHeaderConcern } from '../concern/journal/detail/header';
import { createImageHelper, NewImagesConcern } from '../concern/journal/detail/images';
import { NewImagePreviewConcern } from '../concern/journal/detail/image_preview';
import { NewSidebarConcern } from '../concern/journal/sidebar';
import '../types/platform';

function createJournalDetailPageData(journalId = '') {
	let page = {} as JournalDetailPageData;
	const pg = () => page;
	const image = createImageHelper();

	page.client = NewJournalClient();
	page.noteClient = NewJournalNoteClient();
	page.tagClient = NewJournalTagClient();

	page.presentation = NewPresentationConcern();
	page.current = NewCurrentJournalConcern(pg);
	page.header = NewHeaderConcern(pg);
	page.images = NewImagesConcern(pg, image);
	page.preview = NewImagePreviewConcern(pg);
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

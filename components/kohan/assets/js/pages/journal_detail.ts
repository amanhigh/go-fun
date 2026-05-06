import { NewJournalClient } from '../client/journal';
import { NewJournalNoteClient } from '../client/journal_note';
import { NewJournalTagClient } from '../client/journal_tag';
import type { JournalDetailPageData } from '../types/journal_detail_concern';
import { NewPresentationConcern } from '../concern/journal/common/presentation';
import { NewStateConcern } from '../concern/journal/detail/state';
import { NewHeaderConcern } from '../concern/journal/detail/header';
import { createImageHelper, NewImagesConcern } from '../concern/journal/detail/images';
import { NewModalConcern } from '../concern/journal/detail/modal';
import { NewSidebarConcern } from '../concern/journal/detail/sidebar';
import '../types/platform';

function createJournalDetailPageData() {
	let page = {} as JournalDetailPageData;
	const pg = () => page;
	const image = createImageHelper();

	page.client = NewJournalClient();
	page.noteClient = NewJournalNoteClient();
	page.tagClient = NewJournalTagClient();

	Object.assign(page, NewStateConcern(pg));
	Object.assign(page, NewHeaderConcern(pg));
	Object.assign(page, NewImagesConcern(pg, image));
	Object.assign(page, NewModalConcern(pg));
	page.presentation = NewPresentationConcern();
	page.sidebar = NewSidebarConcern(pg);

	page.init = function init(this: any) {
		page = this;
		this.initDetail();
	};

	return page;
}

document.addEventListener('alpine:init', () => {
	window.Alpine.data('journalDetailPage', createJournalDetailPageData);
});

export {};

import { NewJournalClient } from '../client/journal';
import { NewJournalNoteClient } from '../client/journal_note';
import { NewJournalTagClient } from '../client/journal_tag';
import { createJournalDetailFormatters } from './formatters';
import { createImageHelper } from '../journal_images';
import { createJournalDetailPageState } from './page_state';
import { createJournalDetailPageActions } from './page_actions';
import { createJournalDetailPreview } from './preview_actions';
import { createSideBar } from './side_bar';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalDetailPage>): void;
};

function journalDetailPage() {
	const journalClient = NewJournalClient();
	const noteClient = NewJournalNoteClient();
	const tagClient = NewJournalTagClient();
	const image = createImageHelper();

	const state = createJournalDetailPageState();
	const formatters = createJournalDetailFormatters();
	// FIXME: Review all journal_detail ts files.
	const pageActions = createJournalDetailPageActions({ journalClient });
	const preview = createJournalDetailPreview(image);

	const page: any = {
		...state,
		...formatters,
		...pageActions,
		...preview,
	};

	page.sidebar = createSideBar(page, journalClient, noteClient, tagClient);

	return page;
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalDetailPage', journalDetailPage);
});

export {};

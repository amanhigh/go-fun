import { NewJournalClient } from '../client/journal';
import { NewJournalNoteClient } from '../client/journal_note';
import { NewJournalTagClient } from '../client/journal_tag';
import { createJournalDetailFormatters } from './formatters';
import { createImageHelper } from '../journal_images';
import { createJournalDetailPageState } from './page_state';
import { createJournalDetailPageActions } from './page_actions';
import { createJournalDetailNotes } from './notes_actions';
import { createJournalDetailPreview } from './preview_actions';
import { createJournalDetailReview } from './review_actions';
import { createJournalDetailTags } from './tags_actions';

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
	const pageActions = createJournalDetailPageActions({ journalClient });

	return Object.assign(
		state,
		formatters,
		pageActions,
		createJournalDetailPreview(image),
		createJournalDetailReview(journalClient),
		createJournalDetailNotes(noteClient),
		createJournalDetailTags(tagClient),
	);
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalDetailPage', journalDetailPage);
});

export {};

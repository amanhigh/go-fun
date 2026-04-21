import type { JournalClient } from '../client/journal';
import type { JournalNoteClient } from '../client/journal_note';
import type { JournalTagClient } from '../client/journal_tag';
import { createNotesState } from './notes_state';
import { createJournalDetailNotes } from './notes_actions';
import { createReviewState } from './review_state';
import { createJournalDetailReview } from './review_actions';
import { createTagsState } from './tags_state';
import { createJournalDetailTags } from './tags_actions';

export function createSideBar(
	parent: any,
	journalClient: JournalClient,
	noteClient: JournalNoteClient,
	tagClient: JournalTagClient,
) {
	return Object.assign(
		createNotesState(),
		createReviewState(),
		createTagsState(),
		createJournalDetailNotes(parent, noteClient),
		createJournalDetailReview(parent, journalClient),
		createJournalDetailTags(parent, tagClient),
	);
}

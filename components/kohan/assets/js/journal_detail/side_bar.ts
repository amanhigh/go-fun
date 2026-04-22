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
		{
			get reviewFeedbackClass() {
				const self = this as any;
				return self.reviewMessageType === 'success' ? 'text-emerald-700' : 'text-rose-700';
			},
			get noteFeedbackClass() {
				const self = this as any;
				return self.noteMessageType === 'success' ? 'text-emerald-700' : 'text-rose-700';
			},
			get reasonTagFeedbackClass() {
				const self = this as any;
				return self.reasonTagMessageType === 'success' ? 'text-emerald-700' : 'text-rose-700';
			},
			get managementTagFeedbackClass() {
				const self = this as any;
				return self.managementTagMessageType === 'success' ? 'text-emerald-700' : 'text-rose-700';
			},
			reviewQueueItemClass(value: string) {
				return parent.reviewQueueItemClass(value);
			},
		},
	);
}

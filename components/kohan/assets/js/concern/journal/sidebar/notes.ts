import { createDeletableSyncedCollectionState } from './collection';
import type { JournalNote } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewNotesConcern(pg: JournalDetailPageProvider) {
	return {
		...createDeletableSyncedCollectionState<JournalNote>(
			() => !!pg().current.journal,
			(noteId) => pg().noteClient.delete(pg().current.journalId, noteId),
		),
	};
}

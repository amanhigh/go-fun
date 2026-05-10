import { createDeletableSyncedCollectionState } from '../../../lib/collection';
import type { JournalNote } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function NewNotesConcern(pg: JournalDetailPageProvider) {
	return {
		...createDeletableSyncedCollectionState<JournalNote>(
			() => !!pg().current.journal,
			(noteId) => pg().noteClient.delete(pg().current.journalId, noteId),
		),
		sorted() {
			return [...this.items].sort((left, right) => {
				const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
				const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
				return rightTime - leftTime;
			});
		},
	};
}

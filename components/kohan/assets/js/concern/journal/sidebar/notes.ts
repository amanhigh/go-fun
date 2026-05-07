import { createDeletableSyncedCollectionState } from '../../../lib/collection';
import type { JournalNote } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewNotesConcern(pg: JournalDetailPageProvider) {
	return {
		...createDeletableSyncedCollectionState<JournalNote>(
			() => !!pg().current.journal,
			(noteId) => pg().noteClient.delete(pg().current.journalId, noteId),
			{
				sort: (notes) =>
					[...notes].sort((left, right) => {
						const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
						const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
						return rightTime - leftTime;
					}),
			},
		),
	};
}

import { removeById } from '../../../shared/collection';
import type { JournalNote } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function sortNotes(notes: JournalNote[]): JournalNote[] {
	return [...notes].sort((left: JournalNote, right: JournalNote) => {
		const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
		const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
		return rightTime - leftTime;
	});
}

export function NewNotesConcern(pg: JournalDetailPageProvider) {
	return {
		items: [] as JournalNote[],
		deletingId: '',

		sync(notes: JournalNote[] | undefined) {
			this.items = [...(notes ?? [])];
		},
		sorted() {
			return sortNotes(this.items);
		},
		hasNotes() {
			return this.sorted().length > 0;
		},
		async delete(noteId: string) {
			if (!pg().current.journal) return;
			this.deletingId = noteId;
			try {
				await pg().noteClient.delete(pg().current.journalId, noteId);
				this.items = removeById(this.items, noteId);
			} finally {
				this.deletingId = '';
			}
		},
	};
}

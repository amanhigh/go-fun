import { removeById } from '../../../shared/collection';
import { getErrorMessage } from '../../../shared/error';
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

		sync(this: any, notes: JournalNote[] | undefined) {
			this.items = sortNotes(notes ?? []);
		},
		sorted(this: any) {
			return sortNotes(this.items ?? []);
		},
		async delete(this: any, noteId: string) {
			if (!pg().current.journal || this.deletingId) return;
			this.deletingId = noteId;
			try {
				await pg().noteClient.delete(pg().current.journalId, noteId);
				this.items = sortNotes(removeById(this.items ?? [], noteId));
			} catch (err) {
				getErrorMessage(err, 'Unable to delete note.');
			} finally {
				this.deletingId = '';
			}
		},
	};
}

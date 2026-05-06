import { createAsyncFeedbackState } from '../../../shared/async_feedback';
import { prependById, removeById } from '../../../shared/collection';
import { getErrorMessage } from '../../../shared/error';
import type { JournalNote, JournalNoteRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider, NotesState } from '../../../types/journal_detail_concern';

function sortNotes(notes: JournalNote[]): JournalNote[] {
	return [...notes].sort((left: JournalNote, right: JournalNote) => {
		const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
		const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
		return rightTime - leftTime;
	});
}

export function createNotesState(): NotesState {
	return {
		...createAsyncFeedbackState('noteSubmitting', 'noteMessage', 'noteMessageType'),
		noteDeletingId: '',
		noteContent: '',
		noteItems: [],
	};
}

export function NewNotesConcern(pg: JournalDetailPageProvider) {
	return {
		syncNotes(this: any, notes: JournalNote[] | undefined) {
			this.noteItems = sortNotes(notes ?? []);
		},
		sortedNotes(this: any) {
			return sortNotes(this.noteItems ?? []);
		},
		async submitNote(this: any) {
			const journal = pg().journal;
			if (!journal || this.noteSubmitting) return;
			const content = this.noteContent.trim();
			if (!content) {
				this.noteMessage = 'Note content is required.';
				this.noteMessageType = 'error';
				return;
			}
			this.noteSubmitting = true;
			this.noteMessage = '';
			this.noteMessageType = 'error';
			try {
				const payload: JournalNoteRequest = {
					status: journal.status,
					content,
					format: 'MARKDOWN',
				};
				const envelope = await pg().noteClient.create(pg().journalId, payload);
				this.noteItems = sortNotes(prependById(this.noteItems ?? [], envelope.data));
				this.noteContent = '';
				this.noteMessageType = 'success';
				this.noteMessage = 'Note added.';
			} catch (err) {
				this.noteMessage = getErrorMessage(err, 'Unable to save note.');
				this.noteMessageType = 'error';
			} finally {
				this.noteSubmitting = false;
			}
		},
		async deleteNote(this: any, noteId: string) {
			if (!pg().journal || this.noteDeletingId) return;
			this.noteDeletingId = noteId;
			this.noteMessage = '';
			this.noteMessageType = 'error';
			try {
				await pg().noteClient.delete(pg().journalId, noteId);
				this.noteItems = sortNotes(removeById(this.noteItems ?? [], noteId));
				this.noteMessageType = 'success';
				this.noteMessage = 'Note deleted.';
			} catch (err) {
				this.noteMessage = getErrorMessage(err, 'Unable to delete note.');
				this.noteMessageType = 'error';
			} finally {
				this.noteDeletingId = '';
			}
		},
	};
}

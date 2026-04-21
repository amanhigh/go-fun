import type { JournalNote, JournalNoteClient, JournalNoteRequest } from '../client/journal_note';

export function createJournalDetailNotes(noteClient: JournalNoteClient) {
	return {
		sortedNotes(this: any) {
			return [...(this.journal?.notes ?? [])].sort((left: JournalNote, right: JournalNote) => {
				const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
				const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
				return rightTime - leftTime;
			});
		},
		async submitNote(this: any) {
			if (!this.journal || this.noteSubmitting) return;
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
					status: this.journal.status,
					content,
					format: 'MARKDOWN',
				};
				const envelope = await noteClient.create(this.journalId, payload);
				const notes = this.journal.notes ?? [];
				notes.unshift(envelope.data);
				this.journal.notes = notes;
				this.noteContent = '';
				this.noteMessageType = 'success';
				this.noteMessage = 'Note added.';
			} catch (err) {
				this.noteMessage = err instanceof Error ? err.message : 'Unable to save note.';
				this.noteMessageType = 'error';
			} finally {
				this.noteSubmitting = false;
			}
		},
		async deleteNote(this: any, noteId: string) {
			if (!this.journal || this.noteDeletingId) return;
			this.noteDeletingId = noteId;
			this.noteMessage = '';
			this.noteMessageType = 'error';
			try {
				await noteClient.delete(this.journalId, noteId);
				this.journal.notes = (this.journal.notes ?? []).filter((note: JournalNote) => note.id !== noteId);
				this.noteMessageType = 'success';
				this.noteMessage = 'Note deleted.';
			} catch (err) {
				this.noteMessage = err instanceof Error ? err.message : 'Unable to delete note.';
				this.noteMessageType = 'error';
			} finally {
				this.noteDeletingId = '';
			}
		},
	};
}

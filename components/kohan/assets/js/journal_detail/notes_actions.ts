import type { JournalNote, JournalNoteClient, JournalNoteRequest } from '../client/journal_note';

export function createJournalDetailNotes(parent: any, noteClient: JournalNoteClient) {
	return {
		sortedNotes(this: any) {
			return [...(parent.journal?.notes ?? [])].sort((left: JournalNote, right: JournalNote) => {
				const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
				const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
				return rightTime - leftTime;
			});
		},
		async submitNote(this: any) {
			if (!parent.journal || this.noteSubmitting) return;
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
					status: parent.journal.status,
					content,
					format: 'MARKDOWN',
				};
				const envelope = await noteClient.create(parent.journalId, payload);
				const notes = parent.journal.notes ?? [];
				notes.unshift(envelope.data);
				parent.journal.notes = notes;
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
			if (!parent.journal || this.noteDeletingId) return;
			this.noteDeletingId = noteId;
			this.noteMessage = '';
			this.noteMessageType = 'error';
			try {
				await noteClient.delete(parent.journalId, noteId);
				parent.journal.notes = (parent.journal.notes ?? []).filter((note: JournalNote) => note.id !== noteId);
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

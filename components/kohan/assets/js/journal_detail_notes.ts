import type { Envelope, JournalNote, JournalNoteCreate } from './journal_models';

export function createJournalDetailNotes() {
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
				const payload: JournalNoteCreate = {
					status: this.journal.status,
					content,
					format: 'MARKDOWN',
				};
				const response = await fetch(`/v1/api/journals/${this.journalId}/notes`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(payload),
				});
				if (!response.ok) throw new Error(response.status === 404 ? 'Journal not found' : 'Failed to save note');
				const envelope = (await response.json()) as Envelope<JournalNote>;
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
				const response = await fetch(`/v1/api/journals/${this.journalId}/notes/${noteId}`, {
					method: 'DELETE',
				});
				if (!response.ok) throw new Error(response.status === 404 ? 'Note not found' : 'Failed to delete note');
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

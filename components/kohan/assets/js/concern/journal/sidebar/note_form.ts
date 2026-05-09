import { createSubmitter } from '../../../lib/submitter';
import type { JournalNote, JournalNoteRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewNoteFormConcern(pg: JournalDetailPageProvider) {
	return {
		submitter: createSubmitter(),
		content: '',

		canSubmit() {
			return this.content.trim() !== '';
		},

		async submit() {
			const journal = pg().current.journal;
			if (!journal) return;

			const content = this.content.trim();
			if (!content) {
				this.submitter.setError('Note content is required.');
				return;
			}
			await this.submitter.run(
				() => this.createNote(journal.status, content),
				{ success: 'Note added.', error: 'Unable to save note.' },
			);
		},

		async createNote(status: string, content: string) {
			const payload: JournalNoteRequest = {
				status,
				content,
				format: 'MARKDOWN',
			};
			const envelope = await pg().noteClient.create(pg().current.journalId, payload);
			pg().sidebar.notes.prepend(envelope.data as JournalNote);
			this.content = '';
		},
	};
}

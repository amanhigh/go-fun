import { createSubmitter } from '../../../lib/submitter';
import { JournalNoteFormat } from '../../../types/api/journal/enums';
import type { JournalStatus } from '../../../types/api/journal/enums';
import type { JournalNote } from '../../../types/api/journal/response';
import type { JournalNoteRequest } from '../../../types/api/journal/request';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

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

		async createNote(status: JournalStatus, content: string) {
			const payload: JournalNoteRequest = {
				status,
				content,
				format: JournalNoteFormat.MARKDOWN,
			};
			const envelope = await pg().noteClient.create(pg().current.journalId, payload);
			pg().sidebar.notes.prepend(envelope.data as JournalNote);
			this.content = '';
		},
	};
}

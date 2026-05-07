import { createAsyncFeedback } from '../../../lib/async_feedback';
import type { JournalNote, JournalNoteRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewNoteFormConcern(pg: JournalDetailPageProvider) {
	return {
		...createAsyncFeedback(),
		content: '',

		async submit() {
			const journal = pg().current.journal;
			if (!journal || this.submitting) return;
			const content = this.content.trim();
			if (!content) {
				this.setError('Note content is required.');
				return;
			}
			await this.run(async () => {
				const payload: JournalNoteRequest = {
					status: journal.status,
					content,
					format: 'MARKDOWN',
				};
				const envelope = await pg().noteClient.create(pg().current.journalId, payload);
				pg().sidebar.notes.prepend(envelope.data as JournalNote);
				this.content = '';
			}, 'Note added.', 'Unable to save note.');
		},
	};
}

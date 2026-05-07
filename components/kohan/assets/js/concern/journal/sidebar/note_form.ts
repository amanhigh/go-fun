import { createAsyncFeedbackState } from '../../../lib/async_feedback';
import { getErrorMessage } from '../../../lib/error';
import type { JournalNote, JournalNoteRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewNoteFormConcern(pg: JournalDetailPageProvider) {
	return {
		...createAsyncFeedbackState('submitting', 'message', 'messageType'),
		content: '',

		get feedbackClass(): string {
			return this.messageType === 'success' ? 'journal-feedback-success' : 'journal-feedback-error';
		},

		async submit() {
			const journal = pg().current.journal;
			if (!journal || this.submitting) return;
			const content = this.content.trim();
			if (!content) {
				this.message = 'Note content is required.';
				this.messageType = 'error';
				return;
			}
			this.submitting = true;
			this.message = '';
			this.messageType = 'error';
			try {
				const payload: JournalNoteRequest = {
					status: journal.status,
					content,
					format: 'MARKDOWN',
				};
				const envelope = await pg().noteClient.create(pg().current.journalId, payload);
				pg().sidebar.notes.prepend(envelope.data as JournalNote);
				this.content = '';
				this.messageType = 'success';
				this.message = 'Note added.';
			} catch (err) {
				this.message = getErrorMessage(err, 'Unable to save note.');
				this.messageType = 'error';
			} finally {
				this.submitting = false;
			}
		},
	};
}

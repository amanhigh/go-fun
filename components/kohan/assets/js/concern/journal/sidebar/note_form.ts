import { createAsyncFeedbackState } from '../../../shared/async_feedback';
import { getErrorMessage } from '../../../shared/error';
import { prependById } from '../../../shared/collection';
import type { JournalNote, JournalNoteRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function createNoteFormState() {
	return {
		...createAsyncFeedbackState('submitting', 'message', 'messageType'),
		content: '',
	};
}

export function NewNoteFormConcern(pg: JournalDetailPageProvider) {
	return {
		...createNoteFormState(),

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
				pg().sidebar.notes.items = prependById(pg().sidebar.notes.items ?? [], envelope.data as JournalNote);
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

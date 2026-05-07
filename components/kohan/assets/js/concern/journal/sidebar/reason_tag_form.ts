import { createAsyncFeedbackState } from '../../../shared/async_feedback';
import { getErrorMessage } from '../../../shared/error';
import { prependById } from '../../../shared/collection';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function createReasonTagFormState() {
	return {
		...createAsyncFeedbackState('submitting', 'message', 'messageType'),
		input: '',
		override: '',
	};
}

export function NewReasonTagFormConcern(pg: JournalDetailPageProvider) {
	return {
		...createReasonTagFormState(),

		get feedbackClass(): string {
			return this.messageType === 'success' ? 'journal-feedback-success' : 'journal-feedback-error';
		},

		focusOverride(this: any) {
			(pg() as any).$nextTick?.(() => {
				(pg() as any).$refs?.reasonTagOverride?.focus?.();
			});
		},

		async submit(this: any) {
			if (!pg().current.journal || this.submitting) return;
			const tag = this.input.trim();
			if (!tag) {
				this.message = 'Tag is required.';
				this.messageType = 'error';
				return;
			}
			const override = this.override.trim();
			this.submitting = true;
			this.message = '';
			this.messageType = 'error';
			try {
				const payload: JournalTagRequest = {
					tag,
					type: 'REASON',
					...(override ? { override } : {}),
				};
				const envelope = await pg().tagClient.create(pg().current.journalId, payload);
				pg().sidebar.tags.items = prependById(pg().sidebar.tags.items ?? [], envelope.data as JournalTag);
				this.input = '';
				this.override = '';
				this.messageType = 'success';
				this.message = 'Reason tag added.';
			} catch (err) {
				this.message = getErrorMessage(err, 'Unable to save reason tag.');
				this.messageType = 'error';
			} finally {
				this.submitting = false;
			}
		},
	};
}

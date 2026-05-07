import { getErrorMessage } from './error';

export type AsyncFeedback = {
	submitting: boolean;
	message: string;
	feedbackClass: string;
	setError(message: string): void;
	setSuccess(message: string): void;
	run(action: () => Promise<void>, successMessage: string, errorFallback: string): Promise<void>;
};

export function createAsyncFeedback(): AsyncFeedback {
	return {
		submitting: false,
		message: '',
		feedbackClass: 'journal-feedback-error',

		setError(message: string) {
			this.message = message;
			this.feedbackClass = 'journal-feedback-error';
		},

		setSuccess(message: string) {
			this.message = message;
			this.feedbackClass = 'journal-feedback-success';
		},

		async run(action, successMessage, errorFallback) {
			if (this.submitting) return;
			this.submitting = true;
			this.message = '';
			this.feedbackClass = 'journal-feedback-error';

			try {
				await action();
				this.setSuccess(successMessage);
			} catch (err) {
				this.setError(getErrorMessage(err, errorFallback));
			} finally {
				this.submitting = false;
			}
		},
	};
}

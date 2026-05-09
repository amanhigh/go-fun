

export type FeedbackKind = 'error' | 'success';

const feedbackClassMap: Record<FeedbackKind, string> = {
	error: 'journal-feedback-error',
	success: 'journal-feedback-success',
};

export type Feedback = {
	submitting: boolean;
	message: string;
	feedbackKind: FeedbackKind;
	feedbackClass: string;
	setError(message: string): void;
	setSuccess(message: string): void;
	run(action: () => Promise<void>, successMessage: string, errorFallback: string): Promise<void>;
};

export function createFeedback(): Feedback {
	return {
		submitting: false,
		message: '',
		feedbackKind: 'error',
		feedbackClass: feedbackClassMap.error,

		setError(message: string) {
			this.feedbackKind = 'error';
			this.feedbackClass = feedbackClassMap.error;
			this.message = message;
		},

		setSuccess(message: string) {
			this.feedbackKind = 'success';
			this.feedbackClass = feedbackClassMap.success;
			this.message = message;
		},

		async run(action, successMessage, errorFallback) {
			if (this.submitting) return;
			this.submitting = true;
			this.message = '';
			this.feedbackKind = 'error';
			this.feedbackClass = feedbackClassMap.error;

			try {
				await action();
				this.setSuccess(successMessage);
			} catch (err) {
				this.setError((err as Error).message);
			} finally {
				this.submitting = false;
			}
		},
	};
}

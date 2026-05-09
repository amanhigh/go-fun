import { getErrorMessage } from './error';

// ===== Types =====

export type SubmitMessages = {
	success?: string;
	error: string;
};

// ===== CSS Class Helpers =====

const messageClassMap: Record<string, string> = {
	success: 'journal-feedback-success',
	error: 'journal-feedback-error',
};

// ===== Submitter Type =====

export type Submitter = {
	submitting: boolean;
	message: string;
	messageClass: string;

	hasMessage(): boolean;
	setError(message: string): void;
	run(action: () => Promise<void>, messages: SubmitMessages): Promise<void>;
};

// ===== Factory =====

export function createSubmitter(): Submitter {
	function setSuccess(this: Submitter, message: string) {
		this.messageClass = messageClassMap.success;
		this.message = message;
	}

	return {
		submitting: false,
		message: '',
		messageClass: messageClassMap.error,

		hasMessage(this: Submitter) {
			return this.message !== '';
		},

		setError(this: Submitter, message: string) {
			this.messageClass = messageClassMap.error;
			this.message = message;
		},

		async run(this: Submitter, action: () => Promise<void>, messages: SubmitMessages) {
			if (this.submitting) return;

			this.submitting = true;
			this.message = '';

			try {
				await action();
				setSuccess.call(this, messages.success ?? '');
			} catch (err) {
				this.setError(getErrorMessage(err, messages.error));
			} finally {
				this.submitting = false;
			}
		},
	};
}

import { getErrorMessage } from './error';

// ===== Types =====

export type SubmitMessages = {
	success?: string;
	error: string;
};

// ===== CSS Class Constants =====

const successMessageClass = 'journal-feedback-success';
const errorMessageClass = 'journal-feedback-error';

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
	return {
		submitting: false,
		message: '',
		messageClass: errorMessageClass,

		hasMessage(this: Submitter) {
			return this.message !== '';
		},

		setError(this: Submitter, message: string) {
			this.messageClass = errorMessageClass;
			this.message = message;
		},

		async run(this: Submitter, action: () => Promise<void>, messages: SubmitMessages) {
			// Lifecycle: guard duplicate submit → clear previous message → run action → set success/error → reset submitting
			if (this.submitting) return;

			this.submitting = true;
			this.message = '';

			try {
				await action();
				this.messageClass = successMessageClass;
				this.message = messages.success ?? '';
			} catch (err) {
				this.setError(getErrorMessage(err, messages.error));
			} finally {
				this.submitting = false;
			}
		},
	};
}

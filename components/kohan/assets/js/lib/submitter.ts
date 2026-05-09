import { getErrorMessage } from './error';

// ===== Types =====

export type SubmitStatus = 'idle' | 'submitting' | 'success' | 'error';

export type SubmitMessages = {
	success?: string;
	error: string;
};

// ===== CSS Class Helpers =====

const messageClassMap: Record<Exclude<SubmitStatus, 'idle'>, string> = {
	submitting: '',
	success: 'journal-feedback-success',
	error: 'journal-feedback-error',
};

// ===== Submitter Type =====

export type Submitter = {
	status: SubmitStatus;
	submitting: boolean;
	message: string;
	messageClass: string;

	hasMessage(): boolean;
	hasError(): boolean;
	clear(): void;
	setError(message: string): void;
	setSuccess(message: string): void;
	run(action: () => Promise<void>, messages: SubmitMessages): Promise<void>;
};

// ===== Factory =====

export function createSubmitter(): Submitter {
	return {
		status: 'idle',
		submitting: false,
		message: '',
		messageClass: messageClassMap.error,

		clear(this: Submitter) {
			this.status = 'idle';
			this.message = '';
			this.messageClass = messageClassMap.error;
		},

		hasMessage(this: Submitter) {
			return this.message !== '';
		},

		hasError(this: Submitter) {
			return this.status === 'error' && this.message !== '';
		},

		setError(this: Submitter, message: string) {
			this.status = 'error';
			this.messageClass = messageClassMap.error;
			this.message = message;
		},

		setSuccess(this: Submitter, message: string) {
			this.status = 'success';
			this.messageClass = messageClassMap.success;
			this.message = message;
		},

		async run(this: Submitter, action: () => Promise<void>, messages: SubmitMessages) {
			if (this.submitting) return;

			this.status = 'submitting';
			this.submitting = true;
			this.message = '';

			try {
				await action();
				this.setSuccess(messages.success ?? '');
			} catch (err) {
				this.setError(getErrorMessage(err, messages.error));
			} finally {
				this.submitting = false;
				if (this.status === 'submitting') {
					this.status = 'idle';
				}
			}
		},
	};
}

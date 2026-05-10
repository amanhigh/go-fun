import { createRunnerState, type Runner } from './runner';

// ===== Types =====

export type SubmitMessages = {
	success?: string;
};

// ===== CSS Class Constants =====

const successMessageClass = 'journal-feedback-success';
const errorMessageClass = 'journal-feedback-error';

// ===== Submitter Type =====

export type Submitter = Runner & {
	messageClass: string;

	run(action: () => Promise<void>, messages: SubmitMessages): Promise<boolean>;
};

// ===== Factory =====

export function createSubmitter(): Submitter {
	return {
		...createRunnerState(),
		messageClass: errorMessageClass,

		setError(this: Submitter, message: string) {
			this.error = message;
			this.messageClass = errorMessageClass;
		},

		async run(this: Submitter, action: () => Promise<void>, messages: SubmitMessages): Promise<boolean> {
			const outcome = await this.tryRun(action);
			if (outcome.success) {
				this.error = messages.success ?? '';
				this.messageClass = successMessageClass;
			}
			return outcome.success;
		},
	};
}

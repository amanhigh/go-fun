import { createRunnerState, errorVariant, type Runner } from './runner';

// ===== Types =====

export type SubmitMessages = {
	success?: string;
};

// ===== Submitter Type =====

export type Submitter = Runner & {
	hasError(): boolean;
	run(action: () => Promise<void>, messages: SubmitMessages): Promise<boolean>;
};

// ===== Factory =====

export function createSubmitter(): Submitter {
	return {
		...createRunnerState(),

		hasError(this: Submitter) {
			return this.variant === errorVariant;
		},

		async run(this: Submitter, action: () => Promise<void>, messages: SubmitMessages): Promise<boolean> {
			const outcome = await this.tryRun(action);
			if (outcome.success) {
				this.setSuccess(messages.success ?? '');
			}
			return outcome.success;
		},
	};
}

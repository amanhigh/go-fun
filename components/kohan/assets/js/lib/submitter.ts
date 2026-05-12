import { createRunnerState, type Runner } from './runner';

// ===== Types =====

export type SubmitMessages = {
	success?: string;
};

// ===== Submitter Type =====

export type Submitter = Runner & {
	run(action: () => Promise<void>, messages: SubmitMessages): Promise<boolean>;
};

// ===== Factory =====

export function createSubmitter(): Submitter {
	let dismissTimer: ReturnType<typeof setTimeout> | undefined;

	return {
		...createRunnerState(),

		async run(this: Submitter, action: () => Promise<void>, messages: SubmitMessages): Promise<boolean> {
			// Cancel any previous dismiss timer before starting a new submission.
			if (dismissTimer !== undefined) {
				clearTimeout(dismissTimer);
				dismissTimer = undefined;
			}

			const outcome = await this.tryRun(action);
			if (outcome.success) {
				this.setSuccess(messages.success ?? '');
				// Auto-dismiss success message after 3 seconds.
				// Only clears if the message is unchanged and submitter is not in error state.
				dismissTimer = setTimeout(() => {
					if (!this.hasError() && this.message === (messages.success ?? '')) {
						this.clearMessage();
					}
					dismissTimer = undefined;
				}, 3000);
			}
			return outcome.success;
		},
	};
}

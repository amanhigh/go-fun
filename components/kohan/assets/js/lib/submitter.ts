import { createRunnerState, type Runner } from './runner';

// ===== Submitter Type =====

export type Submitter = Runner & {
	// run executes the action. On success it emits the supplied success
	// message as a transient success notification when a non-null string is
	// provided; pass null to skip the success notification. Validation and
	// caught failures are surfaced automatically as persistent error
	// notifications via setError, so callers never manage inline success/error
	// UI state boxes.
	run(action: () => Promise<void>, successMessage: string | null): Promise<boolean>;
};

// ===== Factory =====

export function createSubmitter(): Submitter {
	const base = createRunnerState();

	return {
		...base,

		// setError surfaces validation and caught failures as a prop-free
		// persistent error notification. The local message state is preserved for the
		// shared execution path.
		setError(this: Submitter, message: string) {
			notify({ message, variant: 'error' });
			base.setError.call(this, message);
		},

		async run(this: Submitter, action: () => Promise<void>, successMessage: string | null): Promise<boolean> {
			const outcome = await this.tryRun(action);
			if (outcome.success && successMessage !== null) {
				notify({ message: successMessage, variant: 'success' });
			}
			return outcome.success;
		},
	};
}

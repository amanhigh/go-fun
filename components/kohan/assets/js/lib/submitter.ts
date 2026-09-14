import { createRunnerState } from './runner';
import { notify } from './notification';

// ===== Submitter Type =====

export interface Submitter {
	busy: boolean;
	isBusy(): boolean;
	setError(message: string): void;

	// run executes the action. On success it emits the supplied success
	// message as a transient success notification when a non-null string is
	// provided; pass null to skip the success notification. Validation and
	// caught failures are surfaced automatically as persistent error
	// notifications via setError, so callers never manage inline success/error
	// UI state boxes.
	run(action: () => Promise<void>, successMessage: string | null): Promise<boolean>;
}

// ===== Factory =====

export function createSubmitter(): Submitter {
	const base = createRunnerState();

	return {
		...base,

		// setError surfaces validation and caught failures as a prop-free
		// persistent error notification.
		setError(message: string) {
			notify({ message, variant: 'error' });
		},

		async run(action: () => Promise<void>, successMessage: string | null): Promise<boolean> {
			const outcome = await base.tryRun.call(this, action);
			if (outcome.status === 'error') {
				this.setError(outcome.message);
				return false;
			}
			if (outcome.status === 'success' && successMessage !== null) {
				notify({ message: successMessage, variant: 'success' });
			}
			return outcome.status === 'success';
		},
	};
}

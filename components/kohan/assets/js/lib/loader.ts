import { createRunnerState, RunOutcomeKind, type Runner } from './runner';
import type { Envelope } from '../types/api/common';

// ===== Loader Type =====

export type Loader = Runner & {
	message: string;
	hasError(): boolean;
	setError(message: string): void;
	clearMessage(): void;

	load<TData>(
		action: () => Promise<Envelope<TData>>,
		onSuccess: (data: TData) => void | Promise<void>,
	): Promise<boolean>;
};

// ===== Factory =====

export function createLoader(): Loader {
	const base = createRunnerState();

	return {
		...base,
		message: '',

		hasError(this: Loader) {
			return this.message !== '';
		},

		setError(this: Loader, message: string) {
			this.message = message;
		},

		clearMessage(this: Loader) {
			this.message = '';
		},

		async load<TData>(
			this: Loader,
			action: () => Promise<Envelope<TData>>,
			onSuccess: (data: TData) => void | Promise<void>,
		): Promise<boolean> {
			const outcome = await base.tryRun.call(this, async () => {
				this.clearMessage();
				const envelope = await action();
				if (envelope.data) {
					await onSuccess(envelope.data);
				}
			});

			if (outcome.kind === RunOutcomeKind.ERROR) {
				this.setError(outcome.error.message);
			}

			return outcome.kind === RunOutcomeKind.SUCCESS;
		},
	};
}

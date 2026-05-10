import { createRunnerState, type Runner } from './runner';
import type { Envelope } from '../types/api/common';

// ===== Loader Type =====

export type Loader = Runner & {
	load<TData>(
		action: () => Promise<Envelope<TData>>,
		onSuccess: (data: TData) => void | Promise<void>,
	): Promise<boolean>;
};

// ===== Factory =====

export function createLoader(): Loader {
	return {
		...createRunnerState(),

		async load<TData>(
			this: Loader,
			action: () => Promise<Envelope<TData>>,
			onSuccess: (data: TData) => void | Promise<void>,
		): Promise<boolean> {
			const outcome = await this.tryRun(async () => {
				const envelope = await action();
				if (envelope.data) {
					await onSuccess(envelope.data);
				}
			});

			return outcome.success;
		},
	};
}

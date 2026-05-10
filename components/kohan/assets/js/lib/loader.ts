import { createRunnerState, type Runner } from './runner';
import type { Envelope } from '../types/api/common';

// ===== Loader Type =====

export type Loader = Runner & {
	loadData<TData>(
		action: () => Promise<Envelope<TData>>,
	): Promise<TData | undefined>;
};

// ===== Factory =====

export function createLoader(): Loader {
	return {
		...createRunnerState(),

		async loadData<TData>(this: Loader, action: () => Promise<Envelope<TData>>): Promise<TData | undefined> {
			const outcome = await this.tryRun(action);
			const envelope = outcome.result;
			return envelope?.data;
		},
	};
}

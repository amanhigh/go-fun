// ===== Types =====

export const RunOutcomeKind = {
	SUCCESS: 'success',
	BUSY: 'busy',
	ERROR: 'error',
} as const;

export type RunOutcomeKind = (typeof RunOutcomeKind)[keyof typeof RunOutcomeKind];

export type RunOutcome<T = void> = {
	kind: typeof RunOutcomeKind.SUCCESS;
	value: T;
} | {
	kind: typeof RunOutcomeKind.BUSY;
} | {
	kind: typeof RunOutcomeKind.ERROR;
	error: Error;
};

export interface Runner {
	busy: boolean;

	isBusy(): boolean;
	tryRun<T>(action: () => Promise<T>): Promise<RunOutcome<T>>;
}

function normalizeError(error: unknown): Error {
	return error instanceof Error ? error : new Error(String(error));
}

// ===== Factory =====

export function createRunnerState(): Runner {
	return {
		busy: false,

		isBusy(this: Runner) {
			return this.busy;
		},

		async tryRun<T>(this: Pick<Runner, 'busy'>, action: () => Promise<T>): Promise<RunOutcome<T>> {
			if (this.busy) return { kind: RunOutcomeKind.BUSY };

			this.busy = true;

			try {
				const result = await action();
				return { kind: RunOutcomeKind.SUCCESS, value: result };
			} catch (err) {
				return {
					kind: RunOutcomeKind.ERROR,
					error: normalizeError(err),
				};
			} finally {
				this.busy = false;
			}
		},
	};
}

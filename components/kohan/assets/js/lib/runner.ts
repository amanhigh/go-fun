// ===== Types =====

export type RunOutcome<T = void> = {
	success: boolean;
	result?: T;
};

export interface Runner {
	busy: boolean;
	error: string;

	hasError(): boolean;
	isBusy(): boolean;
	setError(message: string): void;

	tryRun<T>(action: () => Promise<T>): Promise<RunOutcome<T>>;
}

// ===== Factory =====

export function createRunnerState(): Runner {
	return {
		busy: false,
		error: '',

		hasError(this: Runner) {
			return this.error !== '';
		},

		isBusy(this: Runner) {
			return this.busy;
		},

		setError(this: Runner, message: string) {
			this.error = message;
		},

		async tryRun<T>(this: Runner, action: () => Promise<T>): Promise<RunOutcome<T>> {
			if (this.busy) return { success: false };

			this.busy = true;
			this.error = '';

			try {
				const result = await action();
				return { success: true, result };
			} catch (err) {
				this.setError((err as Error).message);
				return { success: false };
			} finally {
				this.busy = false;
			}
		},
	};
}

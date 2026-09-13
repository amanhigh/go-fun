// ===== Types =====

export type RunOutcome<T = void> = {
	success: boolean;
	result?: T;
};

export interface Runner {
	busy: boolean;
	message: string;

	isBusy(): boolean;
	hasError(): boolean;
	setError(message: string): void;
	clearMessage(): void;

	tryRun<T>(action: () => Promise<T>): Promise<RunOutcome<T>>;
}

// ===== Factory =====

export function createRunnerState(): Runner {
	return {
		busy: false,
		message: '',

		isBusy(this: Runner) {
			return this.busy;
		},

		hasError(this: Runner) {
			return this.message !== '';
		},

		setError(this: Runner, message: string) {
			this.message = message;
		},

		clearMessage(this: Runner) {
			this.message = '';
		},

		async tryRun<T>(this: Runner, action: () => Promise<T>): Promise<RunOutcome<T>> {
			if (this.busy) return { success: false };

			this.busy = true;
			this.clearMessage();

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

// ===== CSS Class Constants =====

export const errorVariant = 'feedback-error';
export const successVariant = 'feedback-success';

// ===== Types =====

export type RunOutcome<T = void> = {
	success: boolean;
	result?: T;
};

export interface Runner {
	busy: boolean;
	message: string;
	variant: string;

	hasMessage(): boolean;
	isBusy(): boolean;
	hasError(): boolean;
	setError(message: string): void;
	setSuccess(message: string): void;
	clearMessage(): void;

	tryRun<T>(action: () => Promise<T>): Promise<RunOutcome<T>>;
}

// ===== Factory =====

export function createRunnerState(): Runner {
	return {
		busy: false,
		message: '',
		variant: '',

		hasMessage(this: Runner) {
			return this.message !== '';
		},

		isBusy(this: Runner) {
			return this.busy;
		},

		hasError(this: Runner) {
			return this.variant === errorVariant;
		},

		setError(this: Runner, message: string) {
			this.message = message;
			this.variant = errorVariant;
		},

		setSuccess(this: Runner, message: string) {
			this.message = message;
			this.variant = successVariant;
		},

		clearMessage(this: Runner) {
			this.message = '';
			this.variant = '';
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

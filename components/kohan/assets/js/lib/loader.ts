// ===== Types =====

export type LoadMessages = {
	error: string;
};

// ===== Loader Type =====

export type Loader = {
	loading: boolean;
	error: string;

	isLoading(): boolean;
	isError(): boolean;
	hasError(): boolean;
	setError(message: string): void;
	run(action: () => Promise<void>, messages: LoadMessages): Promise<boolean>;
};

// ===== Factory =====

export function createLoader(initialLoading = false): Loader {
	return {
		loading: initialLoading,
		error: '',

		isLoading(this: Loader) {
			return this.loading;
		},

		isError(this: Loader) {
			return this.error !== '';
		},

		hasError(this: Loader) {
			return this.error !== '';
		},

		setError(this: Loader, message: string) {
			this.error = message;
		},

		async run(this: Loader, action: () => Promise<void>, messages: LoadMessages): Promise<boolean> {
			// Lifecycle: clear previous error → run action → set error on failure → reset loading
			this.loading = true;
			this.error = '';

			try {
				await action();
				return true;
			} catch (err) {
				this.error = (err as Error).message;
				return false;
			} finally {
				this.loading = false;
			}
		},
	};
}

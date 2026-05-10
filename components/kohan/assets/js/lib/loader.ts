import type { Envelope } from '../types/journal_api';

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

	loadData<TData>(
		action: () => Promise<Envelope<TData>>,
		messages: LoadMessages,
	): Promise<TData | undefined>;
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

		async loadData<TData>(this: Loader, action: () => Promise<Envelope<TData>>, messages: LoadMessages): Promise<TData | undefined> {
			this.loading = true;
			this.error = '';

			try {
				const envelope = await action();
				return envelope.data;
			} catch (err) {
				this.error = (err as Error).message || messages.error;
				return undefined;
			} finally {
				this.loading = false;
			}
		},
	};
}

export type Identifiable = { id: string };

export interface ItemCollection<T> {
	hasItems(): boolean;
}

export interface Collection<T extends Identifiable> extends ItemCollection<T> {
	sync(items: T[]): void;
	all(): T[];
	prepend(item: T): void;
	remove(itemId: string): void;
}

export interface DeletableSyncedCollection<T extends Identifiable> extends Collection<T> {
	deletingId: string;
	delete(itemId: string): Promise<void>;
}

export interface LoadableCollection<T> extends ItemCollection<T> {
	loading: boolean;
	error: string;

	all(): T[];
	isLoading(): boolean;
	isError(): boolean;
	load(): Promise<void>;
}

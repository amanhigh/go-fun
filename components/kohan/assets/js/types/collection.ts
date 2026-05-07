export type Identifiable = { id: string };

export type SyncedCollectionOptions<T> = {
	sort?: (items: T[]) => T[];
};

export interface ItemCollection<T> {
	hasItems(): boolean;
}

export interface SyncedCollection<T extends Identifiable> extends ItemCollection<T> {
	sync(items: T[] | undefined): void;
	all(): T[];
	sorted(): T[];
	prepend(item: T): void;
	remove(itemId: string): void;
}

export interface DeletableSyncedCollection<T extends Identifiable> extends SyncedCollection<T> {
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

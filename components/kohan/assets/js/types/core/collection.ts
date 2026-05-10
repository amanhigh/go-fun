export type Identifiable = { id: string };

export interface Collection<T extends Identifiable> {
	items: T[];
	hasItems(): boolean;
	sync(items: T[]): void;
	all(): T[];
	prepend(item: T): void;
	remove(itemId: string): void;
}

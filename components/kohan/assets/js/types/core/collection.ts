export type Identifiable = { id: string };

export interface Collection<T extends Identifiable> {
	items: T[];
	sync(items: T[]): void;
	all(): T[];
	prepend(item: T): void;
	remove(itemId: string): void;
}

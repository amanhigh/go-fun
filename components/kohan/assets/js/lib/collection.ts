import type { Identifiable } from '../types/core/collection';

export function createCollection<T extends Identifiable>() {
	return {
		items: [] as T[],

		sync(items: T[]) {
			this.items = [...items];
		},
		all() {
			return this.items;
		},
		prepend(item: T) {
			this.items = [item, ...this.items.filter((c) => c.id !== item.id)];
		},
		remove(itemId: string) {
			this.items = this.items.filter((item) => item.id !== itemId);
		},
	};
}

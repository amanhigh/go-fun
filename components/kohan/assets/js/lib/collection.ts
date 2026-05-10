import type { Identifiable } from '../types/core/collection';

// ===== Array Helpers =====

export function prependById<T extends Identifiable>(items: T[], item: T): T[] {
	return [item, ...items.filter((candidate) => candidate.id !== item.id)];
}

export function removeById<T extends Identifiable>(items: T[], itemId: string): T[] {
	return items.filter((item) => item.id !== itemId);
}

// ===== Base Synced Collection =====

export function createCollection<T extends Identifiable>() {
	return {
		items: [] as T[],

		sync(items: T[]) {
			this.items = [...items];
		},
		all() {
			return this.items;
		},
		hasItems() {
			return this.items.length > 0;
		},
		prepend(item: T) {
			this.items = prependById(this.items, item);
		},
		remove(itemId: string) {
			this.items = removeById(this.items, itemId);
		},
	};
}

// ===== Deletable Synced Collection =====

export function createDeletableSyncedCollectionState<T extends Identifiable>(
	canDelete: () => boolean,
	deleteItem: (itemId: string) => Promise<void>,
) {
	return {
		...createCollection<T>(),
		deletingId: '',

		async delete(itemId: string) {
			if (!canDelete()) return;
			this.deletingId = itemId;
			try {
				await deleteItem(itemId);
				this.remove(itemId);
			} finally {
				this.deletingId = '';
			}
		},
	};
}

// ===== Loadable Collection =====

export function createLoadableCollectionState<T>(
	loader: () => Promise<T[]>,
	fallbackMessage: string,
) {
	return {
		items: [] as T[],
		loading: false,
		error: '',

		all() {
			return this.items;
		},
		isLoading() {
			return this.loading;
		},
		isError() {
			return this.error.length > 0;
		},
		hasItems() {
			return this.items.length > 0;
		},

		async load() {
			this.loading = true;
			this.error = '';
			try {
				this.items = await loader();
			} catch (err) {
				this.error = (err as Error).message;
			} finally {
				this.loading = false;
			}
		},
	};
}

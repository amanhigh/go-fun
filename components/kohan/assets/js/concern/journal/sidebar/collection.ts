import { prependById, removeById } from '../../../lib/collection';
import { getErrorMessage } from '../../../lib/error';
import type { Identifiable, SyncedCollectionOptions } from '../../../types/collection';

// ===== Base Synced Collection =====

export function createSyncedCollectionState<T extends Identifiable>(options?: SyncedCollectionOptions<T>) {
	return {
		items: [] as T[],

		sync(items: T[] | undefined) {
			this.items = [...(items ?? [])];
		},
		all() {
			return this.items;
		},
		sorted() {
			return options?.sort?.(this.items) ?? this.items;
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
	options?: SyncedCollectionOptions<T>,
) {
	return {
		...createSyncedCollectionState<T>(options),
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
				this.error = getErrorMessage(err, fallbackMessage);
			} finally {
				this.loading = false;
			}
		},
	};
}

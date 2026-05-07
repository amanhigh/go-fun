import { prependById, removeById } from '../../../shared/collection';
import { getErrorMessage } from '../../../shared/error';

type Identifiable = { id: string };

// ===== Base Synced Collection =====

export function createSyncedCollectionState<T extends Identifiable>() {
	return {
		items: [] as T[],

		sync(items: T[] | undefined) {
			this.items = [...(items ?? [])];
		},
		all() {
			return this.items;
		},
		hasItems() {
			return this.all().length > 0;
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
		...createSyncedCollectionState<T>(),
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

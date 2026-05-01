type StorageValue = string | boolean;

type LocalStorageClient = {
	get(key: string, fallback?: string): string;
	getBool(key: string, fallback: boolean): boolean;
	set(key: string, value: StorageValue): void;
	setBool(key: string, value: boolean): void;
};

function warn(action: string, key: string, err: unknown) {
	console.warn(`localStorage ${action} failed for key "${key}"`, err);
}

function getItem(key: string, fallback = ''): string {
	try {
		const value = window.localStorage.getItem(key);
		return value === null ? fallback : value;
	} catch (err) {
		warn('read', key, err);
		return fallback;
	}
}

function setItem(key: string, value: string): void {
	try {
		window.localStorage.setItem(key, value);
	} catch (err) {
		warn('write', key, err);
	}
}

export function createLocalStorageClient(): LocalStorageClient {
	return {
		get(key, fallback = '') {
			return getItem(key, fallback);
		},
		getBool(key, fallback) {
			return getItem(key, fallback ? 'true' : 'false') === 'true';
		},
		set(key, value) {
			setItem(key, String(value));
		},
		setBool(key, value) {
			setItem(key, value ? 'true' : 'false');
		},
	};
}

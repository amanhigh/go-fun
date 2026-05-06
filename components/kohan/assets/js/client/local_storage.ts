type LocalStorageClient = {
	getBool(key: string, fallback: boolean): boolean;
	setBool(key: string, value: boolean): void;
};

function warn(action: string, key: string, err: unknown) {
	console.warn(`localStorage ${action} failed for key "${key}"`, err);
}

function readString(key: string, fallback = ''): string {
	try {
		const value = window.localStorage.getItem(key);
		return value === null ? fallback : value;
	} catch (err) {
		warn('read', key, err);
		return fallback;
	}
}

function writeString(key: string, value: string): void {
	try {
		window.localStorage.setItem(key, value);
	} catch (err) {
		warn('write', key, err);
	}
}

export function createLocalStorageClient(): LocalStorageClient {
	return {
		getBool(key, fallback) {
			return readString(key, fallback ? 'true' : 'false') === 'true';
		},
		setBool(key, value) {
			writeString(key, value ? 'true' : 'false');
		},
	};
}

export type UrlStateKey = string;

type UrlStateMapping<Key extends UrlStateKey> = {
	fields: readonly Key[];
	queryMap: Partial<Record<Key, string>>;
	reverseMap: Record<string, Key>;
};

type UrlState<Key extends UrlStateKey> = Record<Key, string>;

export function syncUrlToState<Key extends UrlStateKey>(state: UrlState<Key>, mapping: UrlStateMapping<Key>) {
	const params = new URLSearchParams(window.location.search);
	params.forEach((value, key) => {
		const stateKey = mapping.reverseMap[key];
		if (stateKey) {
			state[stateKey] = value;
		}
	});
}

export function syncStateToUrl<Key extends UrlStateKey>(state: UrlState<Key>, mapping: UrlStateMapping<Key>) {
	const params = new URLSearchParams();
	mapping.fields.forEach((key) => {
		const value = state[key];
		if (value !== '') {
			params.set(mapping.queryMap[key] ?? key, value);
		}
	});

	const nextUrl = params.toString() ? `${window.location.pathname}?${params.toString()}` : window.location.pathname;
	window.history.replaceState({}, '', nextUrl);
}

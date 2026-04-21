import type { JournalFilterState } from './filter_state';
import { journalQueryKeyMap, journalReverseQueryKeyMap } from './filter_config';

export function createFilterUrlActions(filter: JournalFilterState) {
	const reverseQueryKeyMap = journalReverseQueryKeyMap as Record<string, string>;
	const trackedFilters = filter as unknown as Record<string, string>;

	return {
		urlToFilter() {
			const params = new URLSearchParams(window.location.search);
			params.forEach((value, key) => {
				const filterKey = reverseQueryKeyMap[key];
				if (filterKey) {
					trackedFilters[filterKey] = value;
				}
			});
		},
		filterToUrl() {
			const params = new URLSearchParams();
			Object.entries(filter.toQueryParams()).forEach(([key, value]) => {
				if (value !== '') params.set(journalQueryKeyMap[key] ?? key, value);
			});
			const nextUrl = params.toString() ? `${window.location.pathname}?${params.toString()}` : window.location.pathname;
			window.history.replaceState({}, '', nextUrl);
		},
	};
}

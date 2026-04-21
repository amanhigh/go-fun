import type { JournalFilterState } from './filter_state';
import { journalFilterKeys } from './filter_state';
import { journalQueryKeyMap, journalReverseQueryKeyMap } from './filter_config';

export function createFilterUrlActions(filter: JournalFilterState) {
	return {
		urlToFilter() {
			const params = new URLSearchParams(window.location.search);
			params.forEach((value, key) => {
				const filterKey = journalReverseQueryKeyMap[key];
				if (filterKey) {
					filter[filterKey] = value;
				}
			});
		},
		filterToUrl() {
			const params = new URLSearchParams();
			journalFilterKeys.forEach((key) => {
				const value = filter[key];
				if (value !== '') params.set(journalQueryKeyMap[key], value);
			});
			const nextUrl = params.toString() ? `${window.location.pathname}?${params.toString()}` : window.location.pathname;
			window.history.replaceState({}, '', nextUrl);
		},
	};
}

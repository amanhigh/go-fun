import type { JournalFilterState } from './filter_state';
import { journalFields } from './filter_state';
import { journalQueryMap, journalReverseMap } from './filter_config';

export function createFilterUrlActions(filter: JournalFilterState) {
	return {
		urlToFilter() {
			const params = new URLSearchParams(window.location.search);
			params.forEach((value, key) => {
				const filterKey = journalReverseMap[key];
				if (filterKey) {
					filter[filterKey] = value;
				}
			});
		},
		filterToUrl() {
			const params = new URLSearchParams();
			journalFields.forEach((key) => {
				const value = filter[key];
				if (value !== '') params.set(journalQueryMap[key], value);
			});
			const nextUrl = params.toString() ? `${window.location.pathname}?${params.toString()}` : window.location.pathname;
			window.history.replaceState({}, '', nextUrl);
		},
	};
}

import { JournalType, JournalStatus, JournalSequence, JournalSortBy, JournalSortOrder, ReviewedFilter } from '../../../types/journal_api';
import type { JournalFilterKey } from '../../../types/journal_api';
import type { JournalFilterUrlConcern, JournalPageProvider } from '../../../types/journal_list_concern';

// Direct mapping from filter state field to URL query key
const journalQueryMap: Record<JournalFilterKey, string> = {
	ticker: 'search',
	type: 'type',
	status: 'status',
	sequence: 'sequence',
	createdAfter: 'created-after',
	createdBefore: 'created-before',
	reviewed: 'reviewed',
	sortBy: 'sort-by',
	sortOrder: 'sort-order',
};

// Reverse mapping from URL query key to filter state field
const journalReverseMap: Record<string, JournalFilterKey> = {};
for (const [field, queryKey] of Object.entries(journalQueryMap)) {
	journalReverseMap[queryKey] = field as JournalFilterKey;
}

export const journalFilterFields: JournalFilterKey[] = ['ticker', 'type', 'status', 'sequence', 'createdAfter', 'createdBefore', 'reviewed', 'sortBy', 'sortOrder'];

// Valid value sets for enum-backed filter fields when reading from URL
const knownFilterValues: Record<string, ReadonlySet<string>> = {
	type: new Set(['', JournalType.TAKEN, JournalType.REJECTED]),
	status: new Set(['', ...Object.values(JournalStatus)]),
	sequence: new Set(['', ...Object.values(JournalSequence)]),
	sortBy: new Set(Object.values(JournalSortBy)),
	sortOrder: new Set(Object.values(JournalSortOrder)),
	reviewed: new Set(Object.values(ReviewedFilter)),
};

/** Read filter state from browser URL query params. */
export function urlToFilter(pg: JournalPageProvider) {
	const filter = pg().filter;
	const params = new URLSearchParams(window.location.search);

	// Read date preset from URL (relative values, not absolute dates)
	const raw = params.get('date');
	filter.datePreset = (raw === 'today' || raw === 'last7' || raw === 'last30') ? raw : '';

	// Read all other filter fields from URL query params
	params.forEach((value, key) => {
		const stateKey = journalReverseMap[key];
		if (!stateKey) return;

		// Validate enum-backed fields against known sets; warn and skip unknown
		const validValues = knownFilterValues[stateKey];
		if (validValues && !validValues.has(value)) {
			console.warn(`Unknown ${stateKey} value from URL:`, value);
			return;
		}

		(filter as Record<string, string>)[stateKey] = value;
	});
}

/** Build a URL query string from the current filter state. */
function buildFilterUrl(pg: JournalPageProvider): string {
	const filter = pg().filter;
	const params = new URLSearchParams();

	// Write date preset if active (instead of absolute createdAfter/createdBefore)
	if (filter.datePreset) {
		params.set('date', filter.datePreset);
	}

	// Write all filter fields, skipping date range when a preset is active
	const skipDates = !!filter.datePreset;
	journalFilterFields.forEach((key) => {
		if (skipDates && (key === 'createdAfter' || key === 'createdBefore')) return;
		const value = filter[key];
		if (value !== '') {
			params.set(journalQueryMap[key], value);
		}
	});

	return params.toString() ? `${window.location.pathname}?${params.toString()}` : window.location.pathname;
}

/** Replace the browser history entry with the current filter URL. */
export function filterToUrl(pg: JournalPageProvider): void {
	const nextUrl = buildFilterUrl(pg);
	window.history.replaceState({}, '', nextUrl);
}

export function NewFilterUrlConcern(pg: JournalPageProvider): JournalFilterUrlConcern {
	return {
		urlToFilter() { urlToFilter(pg); },
		filterToUrl() { filterToUrl(pg); },
	};
}

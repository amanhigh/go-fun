import type { JournalFilterKey } from '../../../types/journal_api';
import type { JournalFilterConcern, JournalFilterUrlConcern, JournalPageProvider } from '../../../types/journal_list_concern';

type FilterConfigEntry = {
	queryKey?: string;
	aliases?: readonly string[];
};

const journalFilterConfig: Record<JournalFilterKey, FilterConfigEntry> = {
	ticker: { queryKey: 'search', aliases: ['ticker'] },
	type: {},
	status: {},
	sequence: {},
	createdAfter: { queryKey: 'created-after' },
	createdBefore: { queryKey: 'created-before' },
	reviewed: {},
	sortBy: { queryKey: 'sort-by' },
	sortOrder: { queryKey: 'sort-order' },
};

export const journalFilterFields: JournalFilterKey[] = ['ticker', 'type', 'status', 'sequence', 'createdAfter', 'createdBefore', 'reviewed', 'sortBy', 'sortOrder'];

const journalQueryMap: Partial<Record<JournalFilterKey, string>> = journalFilterFields.reduce((queryMap, field) => {
	const entry = journalFilterConfig[field];
	if (!entry.queryKey) return queryMap;
	return { ...queryMap, [field]: entry.queryKey };
}, {} as Partial<Record<JournalFilterKey, string>>);

const journalReverseMap: Record<string, JournalFilterKey> = journalFilterFields.reduce((reverseMap, field) => {
	const queryKey = journalQueryMap[field] ?? field;
	const aliases = journalFilterConfig[field].aliases ?? [];

	return {
		...reverseMap,
		[queryKey]: field,
		...aliases.reduce<Record<string, JournalFilterKey>>((aliasMap, alias) => ({ ...aliasMap, [alias]: field }), {}),
	};
}, {} as Record<string, JournalFilterKey>);

function urlToFilterState(filter: JournalFilterConcern) {
	const params = new URLSearchParams(window.location.search);

	// Read date preset from URL (relative values, not absolute dates)
	const raw = params.get('date');
	filter.datePreset = (raw === 'today' || raw === 'last7' || raw === 'last30') ? raw : '';

	// Read all other filter fields from URL query params
	params.forEach((value, key) => {
		const stateKey = journalReverseMap[key];
		if (stateKey) {
			filter[stateKey] = value;
		}
	});
}

function filterStateToUrl(filter: JournalFilterConcern) {
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
			params.set(journalQueryMap[key] ?? key, value);
		}
	});

	const nextUrl = params.toString() ? `${window.location.pathname}?${params.toString()}` : window.location.pathname;
	window.history.replaceState({}, '', nextUrl);
}

export function NewFilterUrlConcern(pg: JournalPageProvider): JournalFilterUrlConcern {
	const filter = pg().filter;
	return {
		urlToFilter() {
			urlToFilterState(filter);
		},
		filterToUrl() {
			filterStateToUrl(filter);
		},
	};
}

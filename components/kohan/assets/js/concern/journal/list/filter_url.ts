import { syncStateToUrl, syncUrlToState } from '../../../shared/url_state';
import type { JournalFilterKey, JournalFilters } from '../../../types/journal_api';
import type { JournalFilterState, JournalFilterUrlState } from '../../../types/journal_list_state';

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

export const journalFilterUrlMapping = {
	fields: journalFilterFields,
	queryMap: journalQueryMap,
	reverseMap: journalReverseMap,
} as const;

function asJournalFilters(filter: JournalFilterState): JournalFilters {
	return filter as unknown as JournalFilters;
}

export function createJournalFilterUrlConcern(filter: JournalFilterState): JournalFilterUrlState {
	return {
		urlToFilter() {
			syncUrlToState(asJournalFilters(filter), journalFilterUrlMapping);
		},
		filterToUrl() {
			syncStateToUrl(asJournalFilters(filter), journalFilterUrlMapping);
		},
	};
}

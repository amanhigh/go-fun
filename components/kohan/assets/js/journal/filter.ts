export const journalFields = ['ticker', 'type', 'status', 'sequence', 'createdAfter', 'createdBefore', 'reviewed', 'sortBy', 'sortOrder'] as const;

export type JournalFilterKey = typeof journalFields[number];

export const journalQueryMap: Record<JournalFilterKey, string> = {
	ticker: 'search',
	createdAfter: 'created-after',
	createdBefore: 'created-before',
	sortBy: 'sort-by',
	sortOrder: 'sort-order',
};

export const journalReverseMap: Record<string, JournalFilterKey> = {
	search: 'ticker',
	ticker: 'ticker',
	type: 'type',
	status: 'status',
	sequence: 'sequence',
	'created-after': 'createdAfter',
	'created-before': 'createdBefore',
	reviewed: 'reviewed',
	'sort-by': 'sortBy',
	'sort-order': 'sortOrder',
};

export type JournalFilters = Record<JournalFilterKey, string>;

type JournalFilterSnapshot = Pick<JournalFilterState, JournalFilterKey>;

export interface JournalFilterState extends JournalFilters {
	hasFilters(): boolean;
	clear(): void;
	toQueryParams(): JournalFilters;
}

const journalDefaults: JournalFilters = {
	ticker: '',
	type: '',
	status: '',
	sequence: '',
	createdAfter: '',
	createdBefore: '',
	reviewed: '',
	sortBy: '',
	sortOrder: '',
};

function snapshotJournalFilters(filter: JournalFilterSnapshot): JournalFilters {
	const params = {} as JournalFilters;
	journalFields.forEach((field) => {
		params[field] = filter[field];
	});
	return params;
}

export function createJournalFilter(): JournalFilterState {
	return {
		...journalDefaults,
		hasFilters(this: JournalFilterState) {
			return journalFields.some((field) => this[field] !== '');
		},
		clear(this: JournalFilterState) {
			journalFields.forEach((field) => {
				this[field] = '';
			});
		},
		toQueryParams(this: JournalFilterState) {
			return snapshotJournalFilters(this);
		},
	};
}

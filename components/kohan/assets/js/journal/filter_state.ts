export const journalFilterKeys = ['ticker', 'type', 'status', 'sequence', 'createdAfter', 'createdBefore', 'reviewed', 'sortBy', 'sortOrder'] as const;

export type JournalFilterKey = typeof journalFilterKeys[number];

export type JournalFilterValues = Record<JournalFilterKey, string>;

export interface JournalFilterState extends JournalFilterValues {
	hasFilters(): boolean;
	clear(): void;
	toQueryParams(): JournalFilterValues;
}

const journalFilterDefaults: JournalFilterValues = {
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

export function createJournalFilterState(): JournalFilterState {
	return {
		...journalFilterDefaults,
		hasFilters(this: JournalFilterState) {
			return journalFilterKeys.some((field) => this[field] !== '');
		},
		clear(this: JournalFilterState) {
			journalFilterKeys.forEach((field) => {
				this[field] = '';
			});
		},
		toQueryParams(this: JournalFilterState) {
			const params = {} as JournalFilterValues;
			journalFilterKeys.forEach((field) => {
				params[field] = this[field];
			});
			return params;
		},
	};
}

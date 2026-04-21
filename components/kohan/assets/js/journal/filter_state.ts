export const journalFields = ['ticker', 'type', 'status', 'sequence', 'createdAfter', 'createdBefore', 'reviewed', 'sortBy', 'sortOrder'] as const;

export type JournalFilterKey = typeof journalFields[number];

export type JournalFilters = Record<JournalFilterKey, string>;

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

export function createJournalFilterState(): JournalFilterState {
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
			const params = {} as JournalFilters;
			journalFields.forEach((field) => {
				params[field] = this[field];
			});
			return params;
		},
	};
}

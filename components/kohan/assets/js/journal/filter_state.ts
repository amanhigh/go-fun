export interface JournalFilterState {
	ticker: string;
	type: string;
	status: string;
	sequence: string;
	createdAfter: string;
	createdBefore: string;
	reviewed: string;
	sortBy: string;
	sortOrder: string;
	hasFilters(): boolean;
	clear(): void;
	toQueryParams(): Record<string, string>;
}

export function createJournalFilterState(): JournalFilterState {
	return {
		ticker: '',
		type: '',
		status: '',
		sequence: '',
		createdAfter: '',
		createdBefore: '',
		reviewed: '',
		sortBy: '',
		sortOrder: '',
		hasFilters() {
			return [this.ticker, this.type, this.status, this.sequence, this.createdAfter, this.createdBefore, this.reviewed, this.sortBy, this.sortOrder].some((value) => value !== '');
		},
		clear() {
			this.ticker = '';
			this.type = '';
			this.status = '';
			this.sequence = '';
			this.createdAfter = '';
			this.createdBefore = '';
			this.reviewed = '';
			this.sortBy = '';
			this.sortOrder = '';
		},
		toQueryParams() {
			return {
				ticker: this.ticker,
				type: this.type,
				status: this.status,
				sequence: this.sequence,
				createdAfter: this.createdAfter,
				createdBefore: this.createdBefore,
				reviewed: this.reviewed,
				sortBy: this.sortBy,
				sortOrder: this.sortOrder,
			};
		},
	};
}

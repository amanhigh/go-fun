export interface PaginationTracker {
	page: number;
	pageSize: number;
	totalItems: number;
	getPage(): number;
	getPageSize(): number;
	getTotalItems(): number;
	getTotalPages(): number;
	hasNext(): boolean;
	hasPrev(): boolean;
	setTotalItems(count: number): void;
	setPageFromOffset(offset: number): void;
	nextPage(): void;
	prevPage(): void;
	resetPage(): void;
}

export function createPaginationTracker(pageSize: number): PaginationTracker {
	return {
		page: 1,
		pageSize,
		totalItems: 0,
		getPage() { return this.page; },
		getPageSize() { return this.pageSize; },
		getTotalItems() { return this.totalItems; },
		getTotalPages() { return Math.max(1, Math.ceil(this.totalItems / this.pageSize)); },
		hasNext() { return this.page < this.getTotalPages(); },
		hasPrev() { return this.page > 1; },
		setTotalItems(count: number) { this.totalItems = count; },
		setPageFromOffset(offset: number) { this.page = Math.floor(offset / this.pageSize) + 1; },
		nextPage() { if (this.hasNext()) this.page += 1; },
		prevPage() { if (this.hasPrev()) this.page -= 1; },
		resetPage() { this.page = 1; },
	};
}

export interface FilterTracker {
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

export function createFilterTracker(): FilterTracker {
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

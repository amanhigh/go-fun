export interface PaginationState {
	page: number;
	pageSize: number;
	totalItems: number;
	getPage(): number;
	getPageSize(): number;
	getOffset(): number;
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

export function createPaginationState(pageSize: number): PaginationState {
	return {
		page: 1,
		pageSize,
		totalItems: 0,
		getPage() { return this.page; },
		getPageSize() { return this.pageSize; },
		getOffset() { return (this.page - 1) * this.pageSize; },
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

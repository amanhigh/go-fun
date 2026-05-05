import type { JournalPageProvider, PaginationConcern } from '../../../types/journal_list_concern';

const defaultPageSize = 10;

export function newPaginationConcern(pg: JournalPageProvider): PaginationConcern {
	return {
		page: 1,
		pageSize: defaultPageSize,
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
		async previousPage() {
			if (!this.hasPrev()) return;
			this.prevPage();
			await pg().table.loadJournals();
		},
		async nextJournalPage() {
			if (!this.hasNext()) return;
			this.nextPage();
			await pg().table.loadJournals();
		},
		summary() { return `Page ${this.getPage()} of ${this.getTotalPages()}`; },
	};
}

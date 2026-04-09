import { JournalClient } from './journal_client';
import type { Journal } from './journal_models';
import { createPaginationTracker } from './journal_state';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalPage>): void;
};

function journalPage() {
	const client = new JournalClient();
	const pagination = createPaginationTracker(10);
	return {
		journals: [] as Journal[],
		pagination,
		loading: false,
		errorMessage: '',
		async loadJournals() {
			this.loading = true;
			this.errorMessage = '';
			try {
				const resp = await client.list(this.pagination.getPage() === 1 ? 0 : (this.pagination.getPage() - 1) * this.pagination.getPageSize(), this.pagination.getPageSize());
				const data = resp.data ?? {};
				this.journals = data.journals ?? [];
				this.pagination.setTotalItems(data.metadata?.total ?? this.journals.length);
				this.pagination.setPageFromOffset(data.metadata?.offset ?? 0);
			} catch {
				this.errorMessage = 'Unable to load journals. Please try again.';
			} finally {
				this.loading = false;
			}
		},
		hasError() {
			return this.errorMessage !== '';
		},
		isEmpty() {
			return this.journals.length === 0;
		},
		formatTimestamp(value: string) {
			if (!value) return '—';
			return new Date(value).toLocaleString();
		},
		async prevPage() {
			if (!this.pagination.hasPrev()) return;
			this.pagination.prevPage();
			await this.loadJournals();
		},
		async nextPage() {
			if (!this.pagination.hasNext()) return;
			this.pagination.nextPage();
			await this.loadJournals();
		},
		init() {
			void this.loadJournals();
		},
	};
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalPage', journalPage);
});

export {};

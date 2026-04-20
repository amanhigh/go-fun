import { journalQueryKeyMap, journalReverseQueryKeyMap } from './journal_models';

const reviewPresetBaseClass = 'border-cyan-200/70 bg-white/80 text-cyan-800 hover:bg-cyan-50';
const reviewPresetAnchorClass = 'border-2 border-amber-200 bg-white/80 text-cyan-800';
const reviewPresetActiveClass = 'border-amber-300 bg-amber-100/90 text-amber-950 hover:bg-amber-100';

function formatReviewPresetLabel(date: Date): string {
	const monthLabels = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
	return `${monthLabels[date.getMonth()]}-${String(date.getFullYear() % 100).padStart(2, '0')}`;
}

export type ReviewPreset = {
	isAnchor: boolean;
	label: string;
	createdAfter: string;
	createdBefore: string;
};

export function createReviewPresets(): ReviewPreset[] {
	const today = new Date();
	const anchorDate = new Date(today.getFullYear(), today.getMonth() - 9, 1);
	return [-2, -1, 0, 1, 2].map((offset) => {
		const monthDate = new Date(anchorDate.getFullYear(), anchorDate.getMonth() + offset, 1);
		const createdAfter = new Date(monthDate.getFullYear(), monthDate.getMonth(), 1);
		const createdBefore = new Date(monthDate.getFullYear(), monthDate.getMonth() + 1, 0);
		return {
			isAnchor: offset === 0,
			label: formatReviewPresetLabel(monthDate),
			createdAfter: `${createdAfter.getFullYear()}-${String(createdAfter.getMonth() + 1).padStart(2, '0')}-${String(createdAfter.getDate()).padStart(2, '0')}`,
			createdBefore: `${createdBefore.getFullYear()}-${String(createdBefore.getMonth() + 1).padStart(2, '0')}-${String(createdBefore.getDate()).padStart(2, '0')}`,
		};
	});
}

export function createJournalListFilterActions(filterTracker: Record<string, string>) {
	const reverseQueryKeyMap = journalReverseQueryKeyMap as Record<string, string>;

	return {
		currentReviewPresetLabel(this: any) {
			return this.activeReviewPreset;
		},
		syncActiveReviewPreset(this: any) {
			const matchingPreset = this.reviewPresets.find((reviewPreset: ReviewPreset) => (
				reviewPreset.createdAfter === this.filterTracker.createdAfter
				&& reviewPreset.createdBefore === this.filterTracker.createdBefore
				&& this.filterTracker.reviewed === 'false'
			));
			this.activeReviewPreset = matchingPreset?.label ?? '';
		},
		clearActiveReviewPreset(this: any) {
			this.activeReviewPreset = '';
		},
		reviewPresetButtonClass(this: any, reviewPreset: ReviewPreset) {
			if (this.activeReviewPreset === reviewPreset.label) {
				return reviewPresetActiveClass;
			}
			return reviewPreset.isAnchor ? reviewPresetAnchorClass : reviewPresetBaseClass;
		},
		toggleTypeFilter(this: any) {
			this.filterTracker.type = this.filterTracker.type === 'TAKEN' ? 'REJECTED' : 'TAKEN';
			this.applyManualFilters();
		},
		urlToFilter() {
			const params = new URLSearchParams(window.location.search);
			params.forEach((value, key) => {
				const filterKey = reverseQueryKeyMap[key];
				if (filterKey) {
					filterTracker[filterKey] = value;
				}
			});
		},
		filterToUrl() {
			const params = new URLSearchParams();
			Object.entries(filterTracker).forEach(([key, value]) => {
				if (value !== '') params.set(journalQueryKeyMap[key] ?? key, value);
			});
			const nextUrl = params.toString() ? `${window.location.pathname}?${params.toString()}` : window.location.pathname;
			window.history.replaceState({}, '', nextUrl);
		},
		toDateInputValue(date: Date) {
			const year = date.getFullYear();
			const month = `${date.getMonth() + 1}`.padStart(2, '0');
			const day = `${date.getDate()}`.padStart(2, '0');
			return `${year}-${month}-${day}`;
		},
		applyCreatedPreset(this: any, preset: string) {
			this.filterTracker.clear();
			this.clearActiveReviewPreset();
			const today = new Date();
			const endDate = this.toDateInputValue(today);
			const daysMap: Record<string, number> = { today: 0, last7: 7, last30: 30 };
			const days = daysMap[preset] ?? 7;
			if (days === 0) {
				this.filterTracker.createdAfter = endDate;
				this.filterTracker.createdBefore = endDate;
				this.applyFilters();
				return;
			}
			const startDate = new Date(today);
			startDate.setDate(today.getDate() - days);
			this.filterTracker.createdAfter = this.toDateInputValue(startDate);
			this.filterTracker.createdBefore = endDate;
			this.applyFilters();
		},
		applyReviewPreset(this: any, reviewPreset: ReviewPreset) {
			this.filterTracker.clear();
			this.filterTracker.createdAfter = reviewPreset.createdAfter;
			this.filterTracker.createdBefore = reviewPreset.createdBefore;
			this.filterTracker.reviewed = 'false';
			this.activeReviewPreset = reviewPreset.label;
			this.applyFilters();
		},
		applyFilters(this: any) {
			this.pagination.resetPage();
			this.filterToUrl();
			void this.loadJournals();
		},
		onCreatedDateChange(this: any) {
			this.clearActiveReviewPreset();
			this.filterTracker.createdBefore = this.filterTracker.createdAfter;
			this.applyFilters();
		},
		applyManualFilters(this: any) {
			this.clearActiveReviewPreset();
			this.applyFilters();
		},
		toggleSort(this: any, field: string) {
			this.clearActiveReviewPreset();
			this.filterTracker.sortOrder = this.filterTracker.sortBy !== field
				? 'asc'
				: this.filterTracker.sortOrder === 'asc' ? 'desc' : 'asc';
			this.filterTracker.sortBy = field;
			this.applyFilters();
		},
		clearFilters(this: any) {
			this.clearActiveReviewPreset();
			this.filterTracker.clear();
			this.applyFilters();
		},
	};
}

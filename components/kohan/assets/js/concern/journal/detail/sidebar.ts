import { createLocalStorageClient } from '../../../client/local_storage';
import { NewNotesConcern, createNotesState } from './notes';
import { NewReviewConcern, createReviewState } from './review_panel';
import { NewTagsConcern, createTagsState, managementTagPresets } from './tags';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewSidebarConcern(pg: JournalDetailPageProvider) {
	const localStorageClient = createLocalStorageClient();

	return Object.assign(
		createNotesState(),
		createReviewState(),
		createTagsState(managementTagPresets),
		NewNotesConcern(pg),
		NewReviewConcern(pg),
		NewTagsConcern(pg),
		{
			actionOpen: true,
			reviewMode: false,
			actionOpenStorageKey: '',
			reviewModeStorageKey: '',
			initSidebarUiState(this: any, actionOpenStorageKey: string, reviewModeStorageKey: string) {
				this.actionOpenStorageKey = actionOpenStorageKey;
				this.reviewModeStorageKey = reviewModeStorageKey;
				this.actionOpen = localStorageClient.getBool(actionOpenStorageKey, true);
				this.reviewMode = localStorageClient.getBool(reviewModeStorageKey, false);
			},
			setActionOpen(this: any, isOpen: boolean) {
				this.actionOpen = isOpen;
				if (this.actionOpenStorageKey) {
					localStorageClient.setBool(this.actionOpenStorageKey, isOpen);
				}
			},
			setReviewMode(this: any, isReviewMode: boolean) {
				this.reviewMode = isReviewMode;
				if (this.reviewModeStorageKey) {
					localStorageClient.setBool(this.reviewModeStorageKey, isReviewMode);
				}
			},
			toggleActionOpen(this: any) {
				this.setActionOpen(!this.actionOpen);
			},
			enterReviewMode(this: any) {
				this.setReviewMode(true);
			},
			exitReviewMode(this: any) {
				this.setReviewMode(false);
			},
			toggleReviewMode(this: any) {
				if (this.reviewMode) {
					this.exitReviewMode();
					return;
				}
				this.enterReviewMode();
			},
			get reviewFeedbackClass() {
				const self = this as any;
				return self.reviewMessageType === 'success' ? 'text-emerald-700' : 'text-rose-700';
			},
			get noteFeedbackClass() {
				const self = this as any;
				return self.noteMessageType === 'success' ? 'text-emerald-700' : 'text-rose-700';
			},
			get reasonTagFeedbackClass() {
				const self = this as any;
				return self.reasonTagMessageType === 'success' ? 'text-emerald-700' : 'text-rose-700';
			},
			get managementTagFeedbackClass() {
				const self = this as any;
				return self.managementTagMessageType === 'success' ? 'text-emerald-700' : 'text-rose-700';
			},
			reviewQueueItemClass(value: string) {
				return pg().reviewQueueItemClass(value);
			},
		},
	);
}

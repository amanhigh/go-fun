import type { JournalClient } from '../client/journal';
import type { JournalNoteClient } from '../client/journal_note';
import type { JournalTagClient } from '../client/journal_tag';
import { createLocalStorageClient } from '../client/local_storage';
import { createNotesState } from './notes_state';
import { createJournalDetailNotes } from './notes_actions';
import { createReviewState } from './review_state';
import { createJournalDetailReview } from './review_actions';
import { createTagsState } from './tags_state';
import { createJournalDetailTags } from './tags_actions';

export function createSideBar(
	parent: any,
	journalClient: JournalClient,
	noteClient: JournalNoteClient,
	tagClient: JournalTagClient,
) {
	const localStorageClient = createLocalStorageClient();

	return Object.assign(
		createNotesState(),
		createReviewState(),
		createTagsState(),
		createJournalDetailNotes(parent, noteClient),
		createJournalDetailReview(parent, journalClient),
		createJournalDetailTags(parent, tagClient),
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
				return parent.reviewQueueItemClass(value);
			},
		},
	);
}

import { createLocalStorageClient } from '../../../client/local_storage';
import { createDeferredAction } from '../../../lib/deferred_action';
import type { Journal } from '../../../types/api/journal/response';
import type { ReviewAdvanceConcern } from '../../../types/journal/sidebar';

const ACTION_OPEN_STORAGE_KEY = 'kohan.journalDetail.sidebar.actionOpen';
const REVIEW_MODE_STORAGE_KEY = 'kohan.journalDetail.reviewMode';

const DEFAULT_ACTION_OPEN = false;
const DEFAULT_REVIEW_OPEN = false;

// newReviewAdvanceConcern provides conditional delayed navigation after a
// journal is marked reviewed in review mode. It composes the deferred action
// primitive directly so Alpine observes the observable active, message, and
// remainingSeconds fields, and relies on the existing full-page navigation
// pattern for the redirect.
function newReviewAdvanceConcern(): ReviewAdvanceConcern {
	const deferred = createDeferredAction();

	return {
		...deferred,
		schedule(this: ReviewAdvanceConcern, next: Journal) {
			this.start(`Advancing to ${next.ticker}…`, () => {
				window.location.href = `/journal/${next.id}`;
			});
		},
	};
}

export function NewSidebarStateConcern() {
	const localStorageClient = createLocalStorageClient();

	return {
		actionOpen: DEFAULT_ACTION_OPEN,
		reviewOpen: DEFAULT_REVIEW_OPEN,
		noteOpen: false,
		reviewAdvance: newReviewAdvanceConcern(),

		restorePersistedSidebarState() {
			this.actionOpen = localStorageClient.getBool(ACTION_OPEN_STORAGE_KEY, DEFAULT_ACTION_OPEN);
			this.reviewOpen = localStorageClient.getBool(REVIEW_MODE_STORAGE_KEY, DEFAULT_REVIEW_OPEN);
		},
		setActionOpen(isOpen: boolean) {
			this.actionOpen = isOpen;
			localStorageClient.setBool(ACTION_OPEN_STORAGE_KEY, isOpen);
		},
		setReviewOpen(isReviewOpen: boolean) {
			this.reviewOpen = isReviewOpen;
			localStorageClient.setBool(REVIEW_MODE_STORAGE_KEY, isReviewOpen);
			if (!isReviewOpen) this.reviewAdvance.cancel();
		},
		setNoteOpen(isOpen: boolean) {
			this.noteOpen = isOpen;
		},
		enterReviewMode() {
			this.setReviewOpen(true);
		},
	};
}

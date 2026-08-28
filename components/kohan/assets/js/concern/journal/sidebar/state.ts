import { createLocalStorageClient } from '../../../client/local_storage';

const ACTION_OPEN_STORAGE_KEY = 'kohan.journalDetail.sidebar.actionOpen';
const REVIEW_MODE_STORAGE_KEY = 'kohan.journalDetail.reviewMode';

const DEFAULT_ACTION_OPEN = false;
const DEFAULT_REVIEW_OPEN = false;

export function NewSidebarStateConcern() {
	const localStorageClient = createLocalStorageClient();

	return {
		actionOpen: DEFAULT_ACTION_OPEN,
		reviewOpen: DEFAULT_REVIEW_OPEN,
		noteOpen: false,

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
		},
		setNoteOpen(isOpen: boolean) {
			this.noteOpen = isOpen;
		},
		enterReviewMode() {
			this.setReviewOpen(true);
		},
	};
}

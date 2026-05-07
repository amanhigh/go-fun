import { createLocalStorageClient } from '../../../client/local_storage';

const ACTION_OPEN_STORAGE_KEY = 'kohan.journalDetail.sidebar.actionOpen';
const REVIEW_MODE_STORAGE_KEY = 'kohan.journalDetail.reviewMode';

export function NewSidebarUiConcern() {
	const localStorageClient = createLocalStorageClient();

	return {
		actionOpen: true,
		reviewMode: false,

		initSidebarUiState() {
			this.actionOpen = localStorageClient.getBool(ACTION_OPEN_STORAGE_KEY, true);
			this.reviewMode = localStorageClient.getBool(REVIEW_MODE_STORAGE_KEY, false);
		},
		setActionOpen(isOpen: boolean) {
			this.actionOpen = isOpen;
			localStorageClient.setBool(ACTION_OPEN_STORAGE_KEY, isOpen);
		},
		setReviewMode(isReviewMode: boolean) {
			this.reviewMode = isReviewMode;
			localStorageClient.setBool(REVIEW_MODE_STORAGE_KEY, isReviewMode);
		},
		enterReviewMode() {
			this.setReviewMode(true);
		},
	};
}

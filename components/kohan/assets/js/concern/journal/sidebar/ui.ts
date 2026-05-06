import { createLocalStorageClient } from '../../../client/local_storage';

const ACTION_OPEN_STORAGE_KEY = 'kohan.journalDetail.sidebar.actionOpen';
const REVIEW_MODE_STORAGE_KEY = 'kohan.journalDetail.reviewMode';

export function NewSidebarUiConcern() {
	const localStorageClient = createLocalStorageClient();

	return {
		actionOpen: true,
		reviewMode: false,

		initSidebarUiState(this: any) {
			this.actionOpen = localStorageClient.getBool(ACTION_OPEN_STORAGE_KEY, true);
			this.reviewMode = localStorageClient.getBool(REVIEW_MODE_STORAGE_KEY, false);
		},
		setActionOpen(this: any, isOpen: boolean) {
			this.actionOpen = isOpen;
		},
		setReviewMode(this: any, isReviewMode: boolean) {
			this.reviewMode = isReviewMode;
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
	};
}

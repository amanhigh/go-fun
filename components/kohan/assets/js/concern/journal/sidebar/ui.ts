import { createLocalStorageClient } from '../../../client/local_storage';

export function NewSidebarUiConcern() {
	const localStorageClient = createLocalStorageClient();

	return {
		actionOpen: true,
		reviewMode: false,

		initSidebarUiState(this: any, actionOpenStorageKey: string, reviewModeStorageKey: string) {
			this.actionOpen = localStorageClient.getBool(actionOpenStorageKey, true);
			this.reviewMode = localStorageClient.getBool(reviewModeStorageKey, false);
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

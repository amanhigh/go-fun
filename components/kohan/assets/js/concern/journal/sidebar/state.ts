import { createLocalStorageClient } from '../../../client/local_storage';
import type { Journal } from '../../../types/api/journal/response';
import type { ReviewAdvanceConcern } from '../../../types/journal/sidebar';
import type { DismissNotification } from '../../../types/notification';

const ACTION_OPEN_STORAGE_KEY = 'kohan.journalDetail.sidebar.actionOpen';
const REVIEW_MODE_STORAGE_KEY = 'kohan.journalDetail.reviewMode';

const DEFAULT_ACTION_OPEN = false;
const DEFAULT_REVIEW_OPEN = false;

// newReviewAdvanceConcern provides conditional delayed navigation after a
// journal is marked reviewed in review mode. The pending advance is owned by
// the shared notification runtime: schedule() emits an actionable success
// notification whose Stay action cancels the advance and whose expiry
// navigates to the next journal. No Alpine-observable state is exposed.
function newReviewAdvanceConcern(): ReviewAdvanceConcern {
	// The dismiss callback for the pending advance notification remains private
	// implementation state; only schedule() and cancel() are exposed.
	let dismiss: DismissNotification | undefined;

	// clearHandle forgets the stored dismiss callback without invoking it. It is
	// used by callbacks the notification runtime has already dismissed (the
	// action runs after runtime dismissal; onExpire runs before runtime cleanup),
	// so the runtime owns the actual dismissal in those paths.
	function clearHandle(): void {
		dismiss = undefined;
	}

	function clearPending(): void {
		dismiss?.();
		dismiss = undefined;
	}

	return {
		schedule(next: Journal) {
			// Cancel any prior pending advance before scheduling a new one.
			clearPending();
			dismiss = notify({
				title: 'Journal reviewed',
				message: `Advancing to ${next.ticker}…`,
				variant: 'success',
				action: {
					label: 'Stay',
					run: clearHandle,
				},
				onExpire: () => {
					clearHandle();
					window.location.href = `/journal/${next.id}`;
				},
			});
		},
		cancel: clearPending,
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

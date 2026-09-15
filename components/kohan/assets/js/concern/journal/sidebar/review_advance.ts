import { notify } from '../../../lib/notification';
import type { DismissNotification } from '../../../types/notification';
import type { Journal } from '../../../types/api/journal/response';

export interface ReviewAdvanceConcern {
	schedule(next: Journal): void;
	cancel(): void;
}

// NewReviewAdvanceConcern owns the conditional delayed navigation after a
// journal is marked reviewed in review mode. Notification expiry owns the
// delay; Stay, dismiss, cancellation, and replacement all prevent navigation.
export function NewReviewAdvanceConcern(): ReviewAdvanceConcern {
	let dismiss: DismissNotification | undefined;

	// The notification runtime has already removed the item when these callbacks
	// run, so only forget the callback without invoking it again.
	function clearHandle(): void {
		dismiss = undefined;
	}

	function clearPending(): void {
		dismiss?.();
		dismiss = undefined;
	}

	return {
		schedule(next: Journal): void {
			clearPending();
			dismiss = notify({
				title: 'Journal reviewed',
				message: `Advancing to ${next.ticker}…`,
				variant: 'success',
				duration: 1500,
				action: {
					label: 'Stay',
					run: clearHandle,
				},
				onExpire: () => {
					clearHandle();
					// FIXME: Open the second-last timeframe automatically on the next review.
					window.location.href = `/journal/${next.id}`;
				},
			});
		},
		cancel: clearPending,
	};
}

import { getErrorMessage } from './error';

export type FeedbackType = 'error' | 'success';

export function createAsyncFeedbackState(): FeedbackState {
	return {
		submitting: false,
		message: '',
		messageType: 'error',
	};
}

export type FeedbackState = {
	submitting: boolean;
	message: string;
	messageType: FeedbackType;
};

/**
 * Wraps an async action with the standard feedback lifecycle:
 * guard on submitting → set up error state → run action → auto-set feedback → error handling → cleanup.
 */
export async function runAsyncFeedback(
	state: FeedbackState,
	action: () => Promise<void>,
	successMessage: string,
	errorFallback: string,
): Promise<void> {
	if (state.submitting) return;
	state.submitting = true;
	state.message = '';
	state.messageType = 'error';
	try {
		await action();
		state.message = successMessage;
		state.messageType = 'success';
	} catch (err) {
		state.message = getErrorMessage(err, errorFallback);
		state.messageType = 'error';
	} finally {
		state.submitting = false;
	}
}

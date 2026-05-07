import { getErrorMessage } from './error';

export type FeedbackType = 'error' | 'success';

export type AsyncFeedbackState<
	SubmittingKey extends string,
	MessageKey extends string,
	MessageTypeKey extends string,
> = Record<SubmittingKey, boolean> & Record<MessageKey, string> & Record<MessageTypeKey, FeedbackType>;

export function createAsyncFeedbackState<
	SubmittingKey extends string,
	MessageKey extends string,
	MessageTypeKey extends string,
>(
	submittingKey: SubmittingKey,
	messageKey: MessageKey,
	messageTypeKey: MessageTypeKey,
): AsyncFeedbackState<SubmittingKey, MessageKey, MessageTypeKey> {
	return {
		[submittingKey]: false,
		[messageKey]: '',
		[messageTypeKey]: 'error',
	} as AsyncFeedbackState<SubmittingKey, MessageKey, MessageTypeKey>;
}

export type FeedbackState = {
	submitting: boolean;
	message: string;
	messageType: FeedbackType;
};

/**
 * Wraps an async action with the standard feedback lifecycle:
 * guard on submitting → set up error state → run action → error handling → cleanup.
 *
 * The action should set `message` and `messageType` on success before returning.
 */
export async function runAsyncFeedback(
	state: FeedbackState,
	action: () => Promise<void>,
	errorFallback: string,
): Promise<void> {
	if (state.submitting) return;
	state.submitting = true;
	state.message = '';
	state.messageType = 'error';
	try {
		await action();
	} catch (err) {
		state.message = getErrorMessage(err, errorFallback);
		state.messageType = 'error';
	} finally {
		state.submitting = false;
	}
}

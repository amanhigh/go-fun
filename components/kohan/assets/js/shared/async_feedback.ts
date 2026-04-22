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

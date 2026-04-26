import type { Journal } from '../client/journal';
import { createAsyncFeedbackState, type FeedbackType } from '../shared/async_feedback';

export type ReviewState = {
	reviewSubmitting: boolean;
	reviewMessage: string;
	reviewMessageType: FeedbackType;
	reviewQueue: Journal[];
	reviewQueueLoading: boolean;
	reviewQueueError: string;
};

export function createReviewState(): ReviewState {
	return {
		...createAsyncFeedbackState('reviewSubmitting', 'reviewMessage', 'reviewMessageType'),
		reviewQueue: [],
		reviewQueueLoading: false,
		reviewQueueError: '',
	};
}

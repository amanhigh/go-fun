import type { Journal } from '../client/journal';

export type ReviewState = {
	reviewSubmitting: boolean;
	reviewMessage: string;
	reviewMessageType: 'error' | 'success';
	reviewQueue: Journal[];
	reviewQueueLoading: boolean;
	reviewQueueError: string;
};

export function createReviewState(): ReviewState {
	return {
		reviewSubmitting: false,
		reviewMessage: '',
		reviewMessageType: 'error',
		reviewQueue: [],
		reviewQueueLoading: false,
		reviewQueueError: '',
	};
}

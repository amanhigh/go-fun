import { createAsyncFeedbackState, type FeedbackType } from '../shared/async_feedback';

export type NotesState = {
	noteSubmitting: boolean;
	noteDeletingId: string;
	noteContent: string;
	noteMessage: string;
	noteMessageType: FeedbackType;
};

export function createNotesState(): NotesState {
	return {
		...createAsyncFeedbackState('noteSubmitting', 'noteMessage', 'noteMessageType'),
		noteDeletingId: '',
		noteContent: '',
	};
}

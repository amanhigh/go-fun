import type { JournalTag } from '../client/journal_tag';
import { managementTagPresets } from './tags_actions';
import { createAsyncFeedbackState, type FeedbackType } from '../shared/async_feedback';

export type TagsState = {
	managementTagPresets: typeof managementTagPresets;
	managementTagSubmitting: boolean;
	managementTagPendingValue: string;
	managementTagMessage: string;
	managementTagMessageType: FeedbackType;
	reasonTagInput: string;
	reasonTagOverride: string;
	reasonTagSubmitting: boolean;
	tagItems: JournalTag[];
	tagDeletingId: string;
	reasonTagMessage: string;
	reasonTagMessageType: FeedbackType;
};

export function createTagsState(): TagsState {
	return {
		...createAsyncFeedbackState('managementTagSubmitting', 'managementTagMessage', 'managementTagMessageType'),
		...createAsyncFeedbackState('reasonTagSubmitting', 'reasonTagMessage', 'reasonTagMessageType'),
		managementTagPresets,
		managementTagPendingValue: '',
		reasonTagInput: '',
		reasonTagOverride: '',
		tagItems: [],
		tagDeletingId: '',
	};
}

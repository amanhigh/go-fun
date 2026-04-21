import { managementTagPresets } from './tags_actions';

export type TagsState = {
	managementTagPresets: typeof managementTagPresets;
	managementTagSubmitting: boolean;
	managementTagPendingValue: string;
	managementTagMessage: string;
	managementTagMessageType: 'error' | 'success';
	reasonTagInput: string;
	reasonTagOverride: string;
	reasonTagSubmitting: boolean;
	tagDeletingId: string;
	reasonTagMessage: string;
	reasonTagMessageType: 'error' | 'success';
};

export function createTagsState(): TagsState {
	return {
		managementTagPresets,
		managementTagSubmitting: false,
		managementTagPendingValue: '',
		managementTagMessage: '',
		managementTagMessageType: 'error',
		reasonTagInput: '',
		reasonTagOverride: '',
		reasonTagSubmitting: false,
		tagDeletingId: '',
		reasonTagMessage: '',
		reasonTagMessageType: 'error',
	};
}

import type { Journal } from '../client/journal';
import { managementTagPresets } from './tags_actions';

export type JournalDetailPageState = {
	journalId: string;
	journal: Journal | null;
	selectedImageIndex: number;
	loading: boolean;
	errorMessage: string;
	reviewSubmitting: boolean;
	noteSubmitting: boolean;
	noteDeletingId: string;
	noteContent: string;
	reviewMessage: string;
	reviewMessageType: 'error' | 'success';
	noteMessage: string;
	noteMessageType: 'error' | 'success';
	reviewQueue: Journal[];
	reviewQueueLoading: boolean;
	reviewQueueError: string;
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

export function createJournalDetailPageState(): JournalDetailPageState {
	return {
		journalId: '',
		journal: null,
		selectedImageIndex: -1,
		loading: false,
		errorMessage: '',
		reviewSubmitting: false,
		noteSubmitting: false,
		noteDeletingId: '',
		noteContent: '',
		reviewMessage: '',
		reviewMessageType: 'error',
		noteMessage: '',
		noteMessageType: 'error',
		reviewQueue: [],
		reviewQueueLoading: false,
		reviewQueueError: '',
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

import type { Journal, JournalNote, JournalTag } from './journal_api';
import type { AlpineContext } from './platform';

export type JournalDetailPageState = {
	journalId: string;
	journal: Journal | null;
	selectedImageIndex: number;
	loading: boolean;
	journalDeleting: boolean;
	errorMessage: string;
};

export type NotesState = {
	noteSubmitting: boolean;
	noteDeletingId: string;
	noteContent: string;
	noteItems: JournalNote[];
	noteMessage: string;
	noteMessageType: 'success' | 'error';
};

export type ReviewState = {
	reviewSubmitting: boolean;
	reviewMessage: string;
	reviewMessageType: 'success' | 'error';
	reviewQueue: Journal[];
	reviewQueueLoading: boolean;
	reviewQueueError: string;
};

export type ManagementTagPreset = {
	value: string;
	label: string;
};

export type TagsState = {
	managementTagPresets: readonly ManagementTagPreset[];
	managementTagSubmitting: boolean;
	managementTagPendingValue: string;
	managementTagMessage: string;
	managementTagMessageType: 'success' | 'error';
	reasonTagInput: string;
	reasonTagOverride: string;
	reasonTagSubmitting: boolean;
	tagItems: JournalTag[];
	tagDeletingId: string;
	reasonTagMessage: string;
	reasonTagMessageType: 'success' | 'error';
};

export type DetailStateBundle = JournalDetailPageState & NotesState & ReviewState & TagsState;

export type DetailAlpineContext = DetailStateBundle & AlpineContext & {
	reviewQueueItemClass(value: string): string;
};

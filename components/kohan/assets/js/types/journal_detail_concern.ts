import type { Journal, JournalImage, JournalNote, JournalTag } from './journal_api';
import type { JournalClient } from '../client/journal';
import type { JournalNoteClient } from '../client/journal_note';
import type { JournalTagClient } from '../client/journal_tag';
import type { JournalPresentationState } from './journal_state';

export type JournalDetailPageProvider = () => JournalDetailPageData;

export type JournalDetailPageData = {
	// Alpine magic properties
	$el: HTMLElement;
	$nextTick: (callback: () => void) => void;
	$refs: Record<string, HTMLElement>;

	// Core state
	journalId: string;
	journal: Journal | null;
	selectedImageIndex: number;
	loading: boolean;
	journalDeleting: boolean;
	errorMessage: string;

	// Clients (on page for concern access via pg())
	client: JournalClient;
	noteClient: JournalNoteClient;
	tagClient: JournalTagClient;

	// Shared presentation
	presentation: JournalPresentationState;

	// Sidebar nested concern
	sidebar: any;

	// Init
	init(): void;

	// Page-level actions
	loadJournal(): Promise<void>;
	deleteJournal(): Promise<void>;
	hasError(): boolean;
	syncSideBarCollections(): void;

	// Preview methods (from modal concern)
	openImagePreview(index: number): void;
	closeImagePreview(): void;
	prevImage(wrap?: boolean): void;
	nextImage(wrap?: boolean): void;
	previewImage(): JournalImage | null;
	previewImageTimeframe(): string;
	previewImageSrc(): string;
	previewImageLabel(): string;
	previewImageCounter(): string;
	hasImagePreview(): boolean;
	canPrevImage(): boolean;
	canNextImage(): boolean;
	sortedImages(): JournalImage[];
	imageCountLabel(): string;

	// Header formatters
	formatTimestamp: (value: string) => string;
	formatDate: (value: string | null | undefined) => string;
	formatReviewQueueDate: (value: string | null | undefined) => string;
	reviewQueueItemClass: (value: string) => string;
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

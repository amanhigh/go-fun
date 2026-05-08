import type { Journal, JournalNote, JournalTag } from './journal_api';
import type { AsyncFeedback } from '../lib/async_feedback';
import type { DeletableSyncedCollection, LoadableCollection } from './collection';
import type { ManagementTagPreset } from '../lib/management_tags';
import type { QuickAction } from './quick_action';

// ===== Sidebar Sub-Concerns =====

export type SidebarStateConcern = {
	actionOpen: boolean;
	reviewOpen: boolean;
	restorePersistedSidebarState(): void;
	setActionOpen(isOpen: boolean): void;
	setReviewOpen(isReviewOpen: boolean): void;
	enterReviewMode(): void;
};

export type ReviewActionsConcern = AsyncFeedback & {
	actions(): QuickAction[];
};

export type ReviewQueueConcern = LoadableCollection<Journal>;

export type NoteFormConcern = AsyncFeedback & {
	content: string;
	submit(): Promise<void>;
};

export type NotesConcern = DeletableSyncedCollection<JournalNote>;

export type TagCollectionConcern = DeletableSyncedCollection<JournalTag> & {
	reason(): JournalTag[];
	directional(): JournalTag[];
	management(): JournalTag[];
};

export type TagFormConcern = AsyncFeedback & {
	input: string;
	override: string;

	submit(): Promise<void>;
};

export type ManagementTagsConcern = AsyncFeedback & {
	presets: readonly ManagementTagPreset[];
	pendingValue: string;

	hasBar(): boolean;
	hasTag(value: string): boolean;
	buttonClass(preset: ManagementTagPreset): string;
	submit(tagValue: string): Promise<void>;
};

export type JournalDetailSidebarConcern = {
	state: SidebarStateConcern;
	reviewActions: ReviewActionsConcern;
	reviewQueue: ReviewQueueConcern;
	noteForm: NoteFormConcern;
	notes: NotesConcern;
	tags: TagCollectionConcern;
	reasonTagForm: TagFormConcern;
	managementTags: ManagementTagsConcern;
};

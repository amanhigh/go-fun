import type { Journal, JournalNote, JournalTag } from './journal_api';
import type { Submitter } from '../lib/submitter';
import type { DeletableSyncedCollection, LoadableCollection } from './collection';
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

export type ReviewActionsConcern = {
	submitter: Submitter;
	actions(): QuickAction[];
};

export type ReviewQueueConcern = LoadableCollection<Journal>;

export type NoteFormConcern = {
	submitter: Submitter;
	content: string;
	canSubmit(): boolean;
	submit(): Promise<void>;
};

export type TagFormConcern = {
	submitter: Submitter;
	input: string;
	override: string;
	canSubmit(): boolean;
	submit(): Promise<void>;
};

export type TakenTagConcern = {
	submitter: Submitter;
	tags: readonly JournalTag[];
	show(): boolean;
	hasTag(value: string): boolean;
	submit(tagValue: string): Promise<void>;
};

export type TagCollectionConcern = DeletableSyncedCollection<JournalTag> & {
	reason(): JournalTag[];
	directional(): JournalTag[];
	management(): JournalTag[];
};

export type NotesConcern = DeletableSyncedCollection<JournalNote> & {
	sorted(): JournalNote[];
};

export type JournalDetailSidebarConcern = {
	state: SidebarStateConcern;
	reviewActions: ReviewActionsConcern;
	reviewQueue: ReviewQueueConcern;
	noteForm: NoteFormConcern;
	notes: NotesConcern;
	tags: TagCollectionConcern;
	reasonTagForm: TagFormConcern;
	takenTag: TakenTagConcern;
};

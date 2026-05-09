import type { Journal, JournalNote, JournalTag } from './journal_api';
import type { Feedback } from '../lib/feedback';
import type { Submitter } from '../lib/submitter';
import type { DeletableSyncedCollection, LoadableCollection } from './collection';
import type { DisplaySpec } from './presentation_concern';
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

export type ReviewActionsConcern = Feedback & {
	actions(): QuickAction[];
};

export type ReviewQueueConcern = LoadableCollection<Journal>;

export type NoteFormConcern = Feedback & {
	content: string;
	submit(): Promise<void>;
};

export type TagFormConcern = {
	submitter: Submitter;
	input: string;
	override: string;
	canSubmit(): boolean;
	submit(): Promise<void>;
};

export type TakenTagConcern = Feedback & {
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

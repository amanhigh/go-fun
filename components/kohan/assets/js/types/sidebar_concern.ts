import type { Journal, JournalNote, JournalTag } from './journal_api';
import type { Feedback } from '../lib/feedback';
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
export type NoteFormConcern = Feedback & {
export type TagFormConcern = Feedback & {
export type TakenTagConcern = Feedback & {
	tags: readonly DisplaySpec[];

	show(): boolean;
	hasTag(value: string): boolean;
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
	takenTag: TakenTagConcern;
};

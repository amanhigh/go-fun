import type { Journal, JournalNote, JournalTag } from '../api/journal/response';
import type { Loader } from '../../lib/loader';
import type { Submitter } from '../../lib/submitter';
import type { Collection } from '../core/collection';
import type { DisplaySpec } from '../core/present';

// ===== Main Concern =====

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

export type QuickAction = {
	id: string;
	isActive(): boolean;
	display: DisplaySpec;
	apply(): Promise<void>;
};

export type ReviewQueueConcern = Collection<Journal> & {
	loader: Loader;
	load(): Promise<void>;
};

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

export type TagCollectionConcern = Collection<JournalTag> & {
	loader: Loader;
	delete(tagId: string): Promise<void>;
	reason(): JournalTag[];
	directional(): JournalTag[];
	management(): JournalTag[];
};

export type NotesConcern = Collection<JournalNote> & {
	loader: Loader;
	delete(noteId: string): Promise<void>;
	sorted(): JournalNote[];
};

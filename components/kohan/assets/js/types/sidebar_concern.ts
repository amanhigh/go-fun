import type { Journal, JournalNote, JournalTag } from './journal_api';
import type { FeedbackType } from '../shared/async_feedback';

// ===== Base Collection Types =====

type Identifiable = { id: string };

export type SyncedCollection<T extends Identifiable> = {
	items: T[];

	sync(items: T[] | undefined): void;
	all(): T[];
	sorted(): T[];
	hasItems(): boolean;
	prepend(item: T): void;
	remove(itemId: string): void;
};

export type DeletableSyncedCollection<T extends Identifiable> = SyncedCollection<T> & {
	deletingId: string;
	delete(itemId: string): Promise<void>;
};

export type LoadableCollection<T> = {
	items: T[];
	loading: boolean;
	error: string;

	isLoading(): boolean;
	isError(): boolean;
	hasItems(): boolean;
	load(): Promise<void>;
};

// ===== Sidebar Sub-Concerns =====

export type ManagementTagPreset = {
	value: string;
	label: string;
};

export type SidebarStateConcern = {
	actionOpen: boolean;
	reviewOpen: boolean;
	restorePersistedSidebarState(): void;
	setActionOpen(isOpen: boolean): void;
	setReviewOpen(isReviewOpen: boolean): void;
	enterReviewMode(): void;
};

export type QuickActionResult = {
	status: string;
	label: string;
	className: string;
};

export type ReviewActionsConcern = {
	submitting: boolean;
	message: string;
	messageType: FeedbackType;
	readonly feedbackClass: string;

	toggleLabel(): string;
	buttonClass(): string;
	quickAction(): QuickActionResult;
	toggle(): Promise<void>;
	applyQuickStatus(): Promise<void>;
};

export type ReviewQueueConcern = LoadableCollection<Journal>;

export type NoteFormConcern = {
	content: string;
	submitting: boolean;
	message: string;
	messageType: 'success' | 'error';
	readonly feedbackClass: string;
	submit(): Promise<void>;
};

export type NotesConcern = DeletableSyncedCollection<JournalNote>;

export type TagCollectionConcern = DeletableSyncedCollection<JournalTag> & {
	reason(): JournalTag[];
	directional(): JournalTag[];
	management(): JournalTag[];
};

export type ReasonTagFormConcern = {
	input: string;
	override: string;
	submitting: boolean;
	message: string;
	messageType: 'success' | 'error';
	readonly feedbackClass: string;

	focusOverride(): void;
	submit(): Promise<void>;
};

export type ManagementTagsConcern = {
	presets: readonly ManagementTagPreset[];
	submitting: boolean;
	pendingValue: string;
	message: string;
	messageType: 'success' | 'error';
	readonly feedbackClass: string;

	hasBar(): boolean;
	hasTag(value: string): boolean;
	buttonClass(value: string): string;
	submit(tagValue: string): Promise<void>;
};

export type JournalDetailSidebarConcern = {
	state: SidebarStateConcern;
	reviewActions: ReviewActionsConcern;
	reviewQueue: ReviewQueueConcern;
	noteForm: NoteFormConcern;
	notes: NotesConcern;
	tags: TagCollectionConcern;
	reasonTagForm: ReasonTagFormConcern;
	managementTags: ManagementTagsConcern;
};

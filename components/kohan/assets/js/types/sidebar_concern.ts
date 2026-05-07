import type { Journal, JournalNote, JournalTag } from './journal_api';
import type { FeedbackType } from '../shared/async_feedback';

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

export type ReviewQueueConcern = {
	items: Journal[];
	loading: boolean;
	error: string;

	hasItems(): boolean;
	load(): Promise<void>;
};

export type NoteFormConcern = {
	content: string;
	submitting: boolean;
	message: string;
	messageType: 'success' | 'error';
	readonly feedbackClass: string;
	submit(): Promise<void>;
};

export type NotesConcern = {
	items: JournalNote[];
	deletingId: string;

	sync(items: JournalNote[] | undefined): void;
	all(): JournalNote[];
	hasItems(): boolean;
	prepend(item: JournalNote): void;
	remove(itemId: string): void;
	sorted(): JournalNote[];
	delete(noteId: string): Promise<void>;
};

export type TagCollectionConcern = {
	items: JournalTag[];
	deletingId: string;

	sync(tags: JournalTag[] | undefined): void;
	all(): JournalTag[];
	hasItems(): boolean;
	prepend(item: JournalTag): void;
	remove(itemId: string): void;
	reason(): JournalTag[];
	directional(): JournalTag[];
	management(): JournalTag[];
	delete(tagId: string): Promise<void>;
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

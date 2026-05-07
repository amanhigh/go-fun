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
	quickAction(): QuickActionResult | null;
	hasQuickAction(): boolean;
	quickLabel(): string;
	quickButtonClass(): string;
	toggle(): Promise<void>;
	applyQuickStatus(): Promise<void>;
};

export type ReviewQueueConcern = {
	items: Journal[];
	loading: boolean;
	error: string;
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
	deleteError: string;
	sync(items: JournalNote[] | undefined): void;
	sorted(): JournalNote[];
	hasNotes(): boolean;
	delete(noteId: string): Promise<void>;
};

export type TagCollectionConcern = {
	items: JournalTag[];
	deletingId: string;
	deleteError: string;

	sync(items: JournalTag[] | undefined): void;
	all(): JournalTag[];
	hasTags(): boolean;
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

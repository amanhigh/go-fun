import type { Journal, JournalNote, JournalTag } from './journal_api';

// ===== Sidebar Sub-Concerns =====

export type ManagementTagPreset = {
	value: string;
	label: string;
};

export type SidebarUiConcern = {
	actionOpen: boolean;
	reviewMode: boolean;
	initSidebarUiState(): void;
	setActionOpen(isOpen: boolean): void;
	setReviewMode(isReviewMode: boolean): void;
	toggleActionOpen(): void;
	enterReviewMode(): void;
	exitReviewMode(): void;
	toggleReviewMode(): void;
};

export type ReviewActionsConcern = {
	submitting: boolean;
	message: string;
	messageType: 'success' | 'error';
	readonly feedbackClass: string;

	toggleLabel(): string;
	buttonClass(): string;
	quickStatus(): string;
	quickLabel(): string;
	hasQuickAction(): boolean;
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
	sync(items: JournalNote[] | undefined): void;
	sorted(): JournalNote[];
	delete(noteId: string): Promise<void>;
};

export type TagCollectionConcern = {
	items: JournalTag[];
	deletingId: string;

	sync(items: JournalTag[] | undefined): void;
	all(): JournalTag[];
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
	ui: SidebarUiConcern;
	reviewActions: ReviewActionsConcern;
	reviewQueue: ReviewQueueConcern;
	noteForm: NoteFormConcern;
	notes: NotesConcern;
	tags: TagCollectionConcern;
	reasonTagForm: ReasonTagFormConcern;
	managementTags: ManagementTagsConcern;
};

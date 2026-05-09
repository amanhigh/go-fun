import type { Journal, JournalImage } from './journal_api';
import type { Submitter } from '../lib/submitter';
import type { JournalClient } from '../client/journal';
import type { JournalNoteClient } from '../client/journal_note';
import type { JournalTagClient } from '../client/journal_tag';
import type { PresentationConcern } from './presentation_concern';
import type { PresentationConcern as PresentConcern } from './present';

// ===== Provider =====

export type JournalDetailPageProvider = () => JournalDetailPageData;

// ===== Page Data (composed from concerns) =====

export type JournalDetailPageData = {
	client: JournalClient;
	noteClient: JournalNoteClient;
	tagClient: JournalTagClient;

	presentation: PresentationConcern;
	present: PresentConcern;
	current: JournalConcern;
	header: JournalHeaderConcern;
	images: JournalImagesConcern;
	preview: PreviewConcern;
	sidebar: JournalDetailSidebarConcern;

	init(): void;
};

// ===== Page-Level Concerns =====

export type JournalConcern = {
	journalId: string;
	journal: Journal | null;
	loading: boolean;
	errorMessage: string;

	loadJournal(): Promise<void>;
	hasError(): boolean;
};

export type JournalHeaderConcern = {
	submitter: Submitter;
	deleteJournal(): Promise<void>;
};

export type JournalImageView = JournalImage & {
	src: string;
	label: string;
};

export type JournalImagesConcern = {
	countLabel(): string;
	sorted(): JournalImageView[];
};

export type PreviewConcern = {
	index: number;
	timeframe(): string;
	src(): string;
	label(): string;
	fileName(): string;
	counter(): string;
	hasPreview(): boolean;
	close(): void;
	prev(): void;
	next(): void;
	wrapPrev(): void;
	wrapNext(): void;
	open(idx: number): void;
};

import type { JournalDetailSidebarConcern } from './sidebar_concern';

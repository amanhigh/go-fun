import type { Journal, JournalImage } from './journal_api';
import type { JournalClient } from '../client/journal';
import type { JournalNoteClient } from '../client/journal_note';
import type { JournalTagClient } from '../client/journal_tag';
import type { PresentationConcern } from './presentation_concern';

// ===== Provider =====

export type JournalDetailPageProvider = () => JournalDetailPageData;

// ===== Page Data (composed from concerns) =====

export type JournalDetailPageData = {
	client: JournalClient;
	noteClient: JournalNoteClient;
	tagClient: JournalTagClient;

	presentation: PresentationConcern;
	current: CurrentJournalConcern;
	header: JournalHeaderConcern;
	images: JournalImagesConcern;
	preview: ImagePreviewConcern;
	sidebar: JournalDetailSidebarConcern;

	init(): void;
};

// ===== Page-Level Concerns =====

export type CurrentJournalConcern = {
	journalId: string;
	journal: Journal | null;
	loading: boolean;
	errorMessage: string;

	loadJournal(): Promise<void>;
	hasError(): boolean;
};

export type JournalHeaderConcern = {
	deleting: boolean;
	deleteJournal(): Promise<void>;
};

export type JournalImagesConcern = {
	resolveImageSrc(fileName: string, createdAt?: string): string;
	timeframeChipClass(timeframe: string): string;
	countLabel(): string;
	sorted(): JournalImage[];
	tileTitle(imageItem: JournalImage): string;
	tileSrc(imageItem: JournalImage): string;
	tileAlt(imageItem: JournalImage): string;
};

export type ImagePreviewConcern = {
	index: number;
	src(): string;
	label(): string;
	counter(): string;
	hasPreview(): boolean;
	close(): void;
	canPrev(): boolean;
	canNext(): boolean;
	prev(wrap?: boolean): void;
	next(wrap?: boolean): void;
	current(): JournalImage | null;
	timeframe(): string;
	open(idx: number): void;
};

import type { JournalDetailSidebarConcern } from './sidebar_concern';

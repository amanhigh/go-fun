import type { JournalDetail, JournalImage } from '../api/journal/response';
import type { JournalTimeframe } from '../api/journal/enums';
import type { Submitter } from '../../lib/submitter';
import type { JournalNoteClient } from '../../client/journal_note';
import type { JournalTagClient } from '../../client/journal_tag';
import type { JournalDetailSidebarConcern } from './sidebar';
import type { JournalPageBase, PageProvider } from './page';

// ===== Main Page Composition =====

export type JournalDetailPage = JournalPageBase & {
	noteClient: JournalNoteClient;
	tagClient: JournalTagClient;

	current: JournalConcern;
	header: JournalHeaderConcern;
	images: JournalImagesConcern;
	preview: PreviewConcern;
	sidebar: JournalDetailSidebarConcern;
};

export type JournalDetailPageProvider = PageProvider<JournalDetailPage>;

// ===== Page Sub-Concerns =====

export type JournalConcern = {
	journalId: string;
	journal: JournalDetail | null;
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
	timeframe(): JournalTimeframe | '';
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

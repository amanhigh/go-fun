import type { JournalStatus, JournalNoteFormat, JournalTagType } from './enums';

export type JournalUpdateRequest = {
	status?: JournalStatus;
	reviewed_at?: string | null;
};

export type JournalNoteRequest = {
	status: JournalStatus;
	content: string;
	format: JournalNoteFormat;
};

export type JournalTagRequest = {
	tag: string;
	type: JournalTagType;
	override?: string;
};

import type { JournalTimeframe, JournalStatus, JournalNoteFormat, JournalTagType, JournalType, JournalSequence } from './enums';
import type { PaginatedResponse } from '../common';

export type JournalImage = {
	id: string;
	timeframe: JournalTimeframe;
	file_name: string;
	created_at: string;
};

export type JournalNote = {
	id: string;
	status: JournalStatus;
	content: string;
	format: JournalNoteFormat;
	created_at: string;
};

export type JournalTag = {
	id: string;
	tag: string;
	type: JournalTagType;
	override?: string;
	created_at: string;
};

export type Journal = {
	id: string;
	ticker: string;
	sequence: JournalSequence;
	type: JournalType;
	status: JournalStatus;
	created_at: string;
	reviewed_at?: string | null;
	images?: JournalImage[];
	tags?: JournalTag[];
	notes?: JournalNote[];
	deleted_at?: string | null;
};

export type JournalList = {
	journals: Journal[];
	metadata: PaginatedResponse;
};

export type JournalUpdate = {
	id: string;
	status: JournalStatus;
	reviewed_at?: string | null;
};

/** Normalized detail view with always-present association arrays. */
export type JournalDetail = Journal & {
	images: JournalImage[];
	tags: JournalTag[];
	notes: JournalNote[];
};

import type { JournalTag } from './journal_api';

export type PresentationConcern = {
	// --- Type ---
	typeBadgeClass(value: string): string;
	typeDisplay(value: string): string;

	// --- Status ---
	normalizeStatus(value: string): string;
	statusBadgeClass(value: string): string;
	statusDisplay(value: string): string;

	// --- Timeframe ---
	timeframeChipClass(value: string): string;

	// --- Sequence ---
	sequenceLabel(sequence: string | null | undefined): string;

	// --- Tag Labels ---
	reasonTagLabel(tag: JournalTag): string;
	directionalTagLabel(tag: JournalTag): string;

	// --- Timestamp / Date ---
	formatTimestamp(value: string | null | undefined): string;
	formatDate(value: string | null | undefined): string;
	formatReviewQueueDate(value: string | null | undefined): string;

	// --- Review Queue ---
	reviewQueueItemClass(value: string): string;

	// --- Feedback ---
	feedbackClass(type: string): string;
};

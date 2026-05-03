export type JournalPresentationState = {
	normalizeStatus(value: string): string;
	statusBadgeClass(value: string): string;
	typeBadgeClass(value: string): string;
	formatTimestamp(value: string | null | undefined): string;
};

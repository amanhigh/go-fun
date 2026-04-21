import type { Journal } from '../client/journal';

export type JournalDetailPageState = {
	journalId: string;
	journal: Journal | null;
	selectedImageIndex: number;
	loading: boolean;
	errorMessage: string;
};

export function createJournalDetailPageState(): JournalDetailPageState {
	return {
		journalId: '',
		journal: null,
		selectedImageIndex: -1,
		loading: false,
		errorMessage: '',
	};
}

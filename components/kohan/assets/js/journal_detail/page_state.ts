import type { Journal } from '../client/journal';

function createEmptyJournal(): Journal {
	return {
		id: '',
		ticker: '',
		sequence: '',
		type: '',
		status: '',
		created_at: '',
		reviewed_at: null,
		images: [],
		tags: [],
		notes: [],
		deleted_at: null,
	};
}

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
		journal: createEmptyJournal(),
		selectedImageIndex: -1,
		loading: true,
		errorMessage: '',
	};
}

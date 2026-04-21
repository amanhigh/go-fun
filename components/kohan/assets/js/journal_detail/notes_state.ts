export type NotesState = {
	noteSubmitting: boolean;
	noteDeletingId: string;
	noteContent: string;
	noteMessage: string;
	noteMessageType: 'error' | 'success';
};

export function createNotesState(): NotesState {
	return {
		noteSubmitting: false,
		noteDeletingId: '',
		noteContent: '',
		noteMessage: '',
		noteMessageType: 'error',
	};
}

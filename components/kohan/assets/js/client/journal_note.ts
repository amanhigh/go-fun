import { BaseClient, type Envelope } from './base';

export type JournalNote = {
	id: string;
	status: string;
	content: string;
	format?: string;
	created_at?: string;
};

export type JournalNoteRequest = {
	status: string;
	content: string;
	format: 'MARKDOWN' | 'PLAINTEXT';
};

export interface JournalNoteClient {
	create(journalId: string, payload: JournalNoteRequest): Promise<Envelope<JournalNote>>;
	delete(journalId: string, noteId: string): Promise<void>;
}

export class JournalNoteClientImpl extends BaseClient implements JournalNoteClient {
	constructor(baseUrl?: string) {
		super(baseUrl);
	}

	async create(journalId: string, payload: JournalNoteRequest): Promise<Envelope<JournalNote>> {
		return this.requestJsonBody<Envelope<JournalNote>>(
			`/journals/${journalId}/notes`,
			'POST',
			payload,
			'Failed to save note',
			'Journal not found',
		);
	}

	async delete(journalId: string, noteId: string): Promise<void> {
		await this.request(
			`/journals/${journalId}/notes/${noteId}`,
			{ method: 'DELETE' },
			'Failed to delete note',
			'Note not found',
		);
	}
}

export function NewJournalNoteClient(): JournalNoteClient {
	return new JournalNoteClientImpl();
}

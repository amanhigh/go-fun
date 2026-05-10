import { BaseClient } from './base';
import type { JournalNote, JournalNoteRequest, Envelope } from '../types/journal_api';

export interface JournalNoteClient {
	create(journalId: string, payload: JournalNoteRequest): Promise<Envelope<JournalNote>>;
	delete(journalId: string, noteId: string): Promise<void>;
}

export class JournalNoteClientImpl extends BaseClient implements JournalNoteClient {
	constructor() {
		super();
	}

	async create(journalId: string, payload: JournalNoteRequest): Promise<Envelope<JournalNote>> {
		return this.requestJson<Envelope<JournalNote>>(`/journals/${journalId}/notes`, 'POST', 'Failed to save note', 'Journal not found', {}, payload);
	}

	async delete(journalId: string, noteId: string): Promise<void> {
		await this.request(`/journals/${journalId}/notes/${noteId}`, { method: 'DELETE' }, 'Failed to delete note', 'Note not found');
	}
}

export function NewJournalNoteClient(): JournalNoteClient {
	return new JournalNoteClientImpl();
}

import { BaseClient, HttpMethod } from './base';
import type { Envelope } from '../types/api/common';
import type { JournalNote } from '../types/api/journal/response';
import type { JournalNoteRequest } from '../types/api/journal/request';

export interface JournalNoteClient {
	create(journalId: string, payload: JournalNoteRequest): Promise<Envelope<JournalNote>>;
	delete(journalId: string, noteId: string): Promise<void>;
}

export class JournalNoteClientImpl extends BaseClient implements JournalNoteClient {
	constructor() {
		super();
	}

	async create(journalId: string, payload: JournalNoteRequest): Promise<Envelope<JournalNote>> {
		return this.requestJson<Envelope<JournalNote>>(`/journals/${journalId}/notes`, HttpMethod.POST, {}, payload);
	}

	async delete(journalId: string, noteId: string): Promise<void> {
		await this.request(`/journals/${journalId}/notes/${noteId}`, { method: HttpMethod.DELETE });
	}
}

export function NewJournalNoteClient(): JournalNoteClient {
	return new JournalNoteClientImpl();
}

import { BaseClient, HttpMethod } from './base';
import type { JournalTag, JournalTagRequest, Envelope } from '../types/journal_api';

export interface JournalTagClient {
	create(journalId: string, payload: JournalTagRequest): Promise<Envelope<JournalTag>>;
	delete(journalId: string, tagId: string): Promise<void>;
}

export class JournalTagClientImpl extends BaseClient implements JournalTagClient {
	constructor() {
		super();
	}

	async create(journalId: string, payload: JournalTagRequest): Promise<Envelope<JournalTag>> {
		return this.requestJson<Envelope<JournalTag>>(`/journals/${journalId}/tags`, HttpMethod.POST, {}, payload);
	}

	async delete(journalId: string, tagId: string): Promise<void> {
		await this.request(`/journals/${journalId}/tags/${tagId}`, { method: HttpMethod.DELETE });
	}
}

export function NewJournalTagClient(): JournalTagClient {
	return new JournalTagClientImpl();
}

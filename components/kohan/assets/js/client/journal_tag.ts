import { BaseClient } from './base';
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
		return this.requestJson<Envelope<JournalTag>>(`/journals/${journalId}/tags`, 'POST', 'Failed to save tag', 'Journal not found', {}, payload);
	}

	async delete(journalId: string, tagId: string): Promise<void> {
		await this.request(`/journals/${journalId}/tags/${tagId}`, { method: 'DELETE' }, 'Failed to delete tag', 'Tag not found');
	}
}

export function NewJournalTagClient(): JournalTagClient {
	return new JournalTagClientImpl();
}

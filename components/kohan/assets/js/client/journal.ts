import { BaseClient, type Envelope, type QueryValue } from './base';
import { journalFields, journalQueryMap, type Journal, type JournalList, type JournalListRequest, type JournalUpdate, type JournalUpdateRequest } from '../types/journal';

export interface JournalClient {
	list(offset: number, limit: number, filters?: JournalListRequest): Promise<Envelope<JournalList>>;
	get(journalId: string): Promise<Envelope<Journal>>;
	updateReview(journalId: string, payload: JournalUpdateRequest): Promise<Envelope<JournalUpdate>>;
	delete(journalId: string): Promise<void>;
}

export class JournalClientImpl extends BaseClient implements JournalClient {
	constructor() {
		super();
	}

	async list(offset: number, limit: number, filters: JournalListRequest = {}): Promise<Envelope<JournalList>> {
		const query: Record<string, QueryValue> = { offset, limit };
		journalFields.forEach((key) => {
			const value = filters[key];
			if (value !== undefined && value !== '') query[journalQueryMap[key] ?? key] = value;
		});
		return this.requestJson<Envelope<JournalList>>('/journals', 'GET', 'Failed to load journals', 'Journal not found', query);
	}

	async get(journalId: string): Promise<Envelope<Journal>> {
		return this.requestJson<Envelope<Journal>>(`/journals/${journalId}`, 'GET', 'Failed to load journal', 'Journal not found');
	}

	async updateReview(journalId: string, payload: JournalUpdateRequest): Promise<Envelope<JournalUpdate>> {
		return this.requestJson<Envelope<JournalUpdate>>(`/journals/${journalId}`, 'PATCH', 'Failed to update journal status', 'Journal not found', {}, payload);
	}

	async delete(journalId: string): Promise<void> {
		await this.request(`/journals/${journalId}`, { method: 'DELETE' }, 'Failed to delete journal', 'Journal not found');
	}

}

export function NewJournalClient(): JournalClient {
	return new JournalClientImpl();
}

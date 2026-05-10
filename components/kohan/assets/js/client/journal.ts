import { BaseClient, HttpMethod, type QueryValue } from './base';
import type { Envelope } from '../types/api/common';
import type { Journal, JournalList, JournalUpdate } from '../types/api/journal/response';
import type { JournalFilterKey, JournalListRequest, JournalUpdateRequest } from '../types/api/journal/request';

const journalApiFields: JournalFilterKey[] = ['ticker', 'type', 'status', 'sequence', 'createdAfter', 'createdBefore', 'reviewed', 'sortBy', 'sortOrder'];

const journalApiQueryMap: Partial<Record<JournalFilterKey, string>> = {
	ticker: 'search',
	createdAfter: 'created-after',
	createdBefore: 'created-before',
	sortBy: 'sort-by',
	sortOrder: 'sort-order',
};

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
		journalApiFields.forEach((key) => {
			const value = filters[key];
			if (value !== undefined && value !== '') query[journalApiQueryMap[key] ?? key] = value;
		});
		return this.requestJson<Envelope<JournalList>>('/journals', HttpMethod.GET, query);
	}

	async get(journalId: string): Promise<Envelope<Journal>> {
		return this.requestJson<Envelope<Journal>>(`/journals/${journalId}`, HttpMethod.GET);
	}

	async updateReview(journalId: string, payload: JournalUpdateRequest): Promise<Envelope<JournalUpdate>> {
		return this.requestJson<Envelope<JournalUpdate>>(`/journals/${journalId}`, HttpMethod.PATCH, {}, payload);
	}

	async delete(journalId: string): Promise<void> {
		await this.request(`/journals/${journalId}`, { method: HttpMethod.DELETE });
	}

}

export function NewJournalClient(): JournalClient {
	return new JournalClientImpl();
}

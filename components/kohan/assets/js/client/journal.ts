import { BaseClient, type Envelope, type QueryValue } from './base';
import type { JournalImage } from './journal_image';
import type { JournalNote } from './journal_note';
import type { JournalTag } from './journal_tag';
import type { JournalFilterState, JournalFilters } from '../journal/filter_state';
import { journalFields } from '../journal/filter_state';
import { journalQueryMap } from '../journal/filter_config';

export type Journal = {
	id: string;
	ticker: string;
	sequence: string;
	type: string;
	status: string;
	created_at: string;
	reviewed_at?: string | null;
	images?: JournalImage[];
	tags?: JournalTag[];
	notes?: JournalNote[];
	deleted_at?: string | null;
};

export type JournalList = {
	journals?: Journal[];
	metadata?: {
		total?: number;
		offset?: number;
		limit?: number;
	};
};

export type JournalUpdateRequest = {
	status?: string;
	reviewed_at: string | null;
};

export type JournalUpdate = {
	id: string;
	status: string;
	reviewed_at: string | null;
};

export type JournalListRequest = ReturnType<JournalFilterState['toQueryParams']>;

export interface JournalClient {
	list(offset: number, limit: number, filters?: JournalListRequest): Promise<Envelope<JournalList>>;
	get(journalId: string): Promise<Envelope<Journal>>;
	updateReview(journalId: string, payload: JournalUpdateRequest): Promise<Envelope<JournalUpdate>>;
}

export class JournalClientImpl extends BaseClient implements JournalClient {
	constructor() {
		super();
	}

	async list(offset: number, limit: number, filters: JournalListRequest = {}): Promise<Envelope<JournalList>> {
		const query: Record<string, QueryValue> = { offset, limit };
		journalFields.forEach((key) => {
			const value = filters[key as keyof JournalFilters];
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

}

export function NewJournalClient(): JournalClient {
	return new JournalClientImpl();
}

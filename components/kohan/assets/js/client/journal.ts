import { BaseClient, type Envelope, type QueryValue } from './base';
import type { JournalImage } from './journal_image';
import type { JournalNote } from './journal_note';
import type { JournalTag } from './journal_tag';
import { journalQueryKeyMap } from '../journal/filter_config';

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

export interface JournalClient {
	list(offset: number, limit: number, filters?: JournalListRequest): Promise<Envelope<JournalList>>;
	get(journalId: string): Promise<Envelope<Journal>>;
	updateReview(journalId: string, payload: JournalUpdateRequest): Promise<Envelope<JournalUpdate>>;
}

export class JournalClientImpl extends BaseClient implements JournalClient {
	async list(offset: number, limit: number, filters: JournalListRequest = {}): Promise<Envelope<JournalList>> {
		const query: Record<string, QueryValue> = { offset, limit };
		Object.entries(filters).forEach(([key, value]) => {
			if (value !== undefined && value !== '') query[journalQueryKeyMap[key] ?? key] = value;
		});
		return this.requestJson<Envelope<JournalList>>('/journals', {}, 'Failed to load journals', 'Journal not found', query);
	}

	async get(journalId: string): Promise<Envelope<Journal>> {
		return this.requestJson<Envelope<Journal>>(`/journals/${journalId}`, {}, 'Failed to load journal', 'Journal not found');
	}

	async updateReview(journalId: string, payload: JournalUpdateRequest): Promise<Envelope<JournalUpdate>> {
		return this.requestJson<Envelope<JournalUpdate>>(
			`/journals/${journalId}`,
			{ method: 'PATCH', headers: this.jsonHeaders(), body: JSON.stringify(payload) },
			'Failed to update journal status',
			'Journal not found',
		);
	}

}

export function NewJournalClient(): JournalClient {
	return new JournalClientImpl();
}

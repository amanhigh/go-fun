import { BaseClient, HttpMethod } from './base';

export interface JournalImageClient {
	delete(journalId: string, imageId: string): Promise<void>;
}

export class JournalImageClientImpl extends BaseClient implements JournalImageClient {
	constructor() {
		super();
	}

	async delete(journalId: string, imageId: string): Promise<void> {
		await this.request(`/journals/${journalId}/images/${imageId}`, { method: HttpMethod.DELETE });
	}
}

export function NewJournalImageClient(): JournalImageClient {
	return new JournalImageClientImpl();
}

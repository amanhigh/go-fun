import type { Envelope, JournalList } from './journal_models';

export class JournalClient {
	async list(offset: number, limit: number): Promise<Envelope<JournalList>> {
		const params = new URLSearchParams({ offset: String(offset), limit: String(limit) });
		const response = await fetch(`/v1/api/journals?${params.toString()}`);
		if (!response.ok) throw new Error('Failed to load journals');
		return response.json() as Promise<Envelope<JournalList>>;
	}
}

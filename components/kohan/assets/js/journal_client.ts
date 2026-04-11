import { journalQueryKeyMap, type Envelope, type JournalList, type JournalListFilters } from './journal_models';

export class JournalClient {
	async list(offset: number, limit: number, filters: JournalListFilters = {}): Promise<Envelope<JournalList>> {
		const params = new URLSearchParams({ offset: String(offset), limit: String(limit) });
		Object.entries(filters).forEach(([key, value]) => {
			if (value !== undefined && value !== '') params.set(journalQueryKeyMap[key] ?? key, value);
		});
		const response = await fetch(`/v1/api/journals?${params.toString()}`);
		if (!response.ok) throw new Error('Failed to load journals');
		return response.json() as Promise<Envelope<JournalList>>;
	}
}

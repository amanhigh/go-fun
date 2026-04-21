export type QueryValue = string | number | boolean | null | undefined;

export type Envelope<T> = {
	status: string;
	data: T;
};

export abstract class BaseClient {
	protected constructor(protected readonly baseUrl = '/v1/api') {}

	protected jsonHeaders(): Record<string, string> {
		return { 'Content-Type': 'application/json' };
	}

	protected buildUrl(path: string, query: Record<string, QueryValue> = {}): string {
		const searchParams = new URLSearchParams();
		Object.entries(query).forEach(([key, value]) => {
			if (value !== undefined && value !== null && value !== '') {
				searchParams.set(key, String(value));
			}
		});
		const queryString = searchParams.toString();
		return queryString ? `${this.baseUrl}${path}?${queryString}` : `${this.baseUrl}${path}`;
	}

	protected async request(path: string, init: RequestInit = {}, errorMessage: string, notFoundMessage = errorMessage, query: Record<string, QueryValue> = {}): Promise<Response> {
		const response = await fetch(this.buildUrl(path, query), init);
		if (!response.ok) {
			throw new Error(response.status === 404 ? notFoundMessage : errorMessage);
		}
		return response;
	}

	protected async requestJson<T>(path: string, init: RequestInit = {}, errorMessage: string, notFoundMessage = errorMessage, query: Record<string, QueryValue> = {}): Promise<T> {
		const response = await this.request(path, init, errorMessage, notFoundMessage, query);
		return response.json() as Promise<T>;
	}
}

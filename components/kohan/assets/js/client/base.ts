export type QueryValue = string | number | boolean | null | undefined;

export type Envelope<T> = {
	status: string;
	data: T;
};

export abstract class BaseClient {
	protected constructor(protected readonly baseUrl = '/v1/api') {}

	protected buildUrl(path: string, query: Record<string, QueryValue> = {}): string {
		const searchParams = new URLSearchParams(
			Object.entries(query)
				.filter(([, value]) => value !== undefined && value !== null && value !== '')
				.map(([key, value]) => [key, String(value)]),
		);
		const queryString = searchParams.toString();
		return queryString ? `${this.baseUrl}${path}?${queryString}` : `${this.baseUrl}${path}`;
	}

	protected requestJsonBody<T>(
		path: string,
		method: string,
		payload: unknown,
		errorMessage: string,
		notFoundMessage = errorMessage,
		query: Record<string, QueryValue> = {},
		init: RequestInit = {},
	): Promise<T> {
		const headers = new Headers(init.headers);
		headers.set('Content-Type', 'application/json');

		return this.requestJson<T>(
			path,
			{ ...init, method, headers, body: JSON.stringify(payload) },
			errorMessage,
			notFoundMessage,
			query,
		);
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

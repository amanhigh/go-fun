export type QueryValue = string | number | boolean | null | undefined;

export type Envelope<T> = {
	status: string;
	data: T;
};

const apiBaseUrl = '/v1/api';

export abstract class BaseClient {
	protected constructor(protected readonly baseUrl = apiBaseUrl) {}

	protected buildUrl(path: string, query: Record<string, QueryValue> = {}): string {
		const searchParams = new URLSearchParams(
			Object.entries(query)
				.filter(([, value]) => value !== undefined && value !== null && value !== '')
				.map(([key, value]) => [key, String(value)]),
		);
		const queryString = searchParams.toString();
		return queryString ? `${this.baseUrl}${path}?${queryString}` : `${this.baseUrl}${path}`;
	}

	protected async requestJson<T>(
		path: string,
		init: RequestInit = {},
		errorMessage: string,
		notFoundMessage = errorMessage,
		query: Record<string, QueryValue> = {},
		payload?: unknown,
	): Promise<T> {
		const requestInit = payload === undefined
			? init
			: {
				...init,
				headers: new Headers(init.headers),
				body: JSON.stringify(payload),
			};

		if (payload !== undefined) {
			(requestInit.headers as Headers).set('Content-Type', 'application/json');
		}

		const response = await this.request(path, requestInit, errorMessage, notFoundMessage, query);
		return response.json() as Promise<T>;
	}

	private async request(path: string, init: RequestInit = {}, errorMessage: string, notFoundMessage = errorMessage, query: Record<string, QueryValue> = {}): Promise<Response> {
		const response = await fetch(this.buildUrl(path, query), init);
		if (!response.ok) {
			throw new Error(response.status === 404 ? notFoundMessage : errorMessage);
		}
		return response;
	}
}

import type { EnvelopeErrorBody, EnvelopeFailBody } from '../types/journal_api';

export type QueryValue = string | number | boolean | null | undefined;

export const HttpMethod = {
	GET: 'GET',
	POST: 'POST',
	PATCH: 'PATCH',
	DELETE: 'DELETE',
} as const;

export type HttpMethod = (typeof HttpMethod)[keyof typeof HttpMethod];

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
		method: HttpMethod,
		query?: Record<string, QueryValue>,
		payload?: unknown,
	): Promise<T> {
		const requestInit = payload === undefined
			? { method }
			: { method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) };

		const response = await this.request(
			path,
			requestInit,
			query ?? {},
		);
		return response.json() as Promise<T>;
	}

	protected async request(path: string, init: RequestInit = {}, query: Record<string, QueryValue> = {}): Promise<Response> {
		const url = this.buildUrl(path, query);
		const response = await fetch(url, init);
		if (!response.ok) {
			throw new Error(await this.fallbackResponse(response, url));
		}
		return response;
	}

	private async fallbackResponse(response: Response, url: string): Promise<string> {
		try {
			const body = await response.json() as Record<string, unknown>;
			const status = body.status as string;

			if (status === 'fail') {
				const data = body.data as EnvelopeFailBody | undefined;
				if (data?.message) return data.message;
				if (data) return Object.values(data).join('; ');
			}

			if (status === 'error') {
				const errorBody = body as unknown as EnvelopeErrorBody;
				if (errorBody.message) return errorBody.message;
			}

			if (body.message) return body.message as string;
		} catch {
			// response body is not JSON; fall through to fallback
		}

		return response.status === 404 ? 'Not found' : `Request failed: ${url}`;
	}
}

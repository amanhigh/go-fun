export type EnvelopeStatus = 'success' | 'fail' | 'error';

export type PaginatedResponse = {
	total: number;
	offset: number;
	limit: number;
};

export type Envelope<T> = {
	status: EnvelopeStatus;
	data: T;
};

// 4xx: { status: "fail", data: { "message"|field: "..." } }
export type EnvelopeFailBody = Record<string, string>;

// 5xx: { status: "error", message: "...", code: ... }
export type EnvelopeErrorBody = {
	message: string;
	code: number;
};

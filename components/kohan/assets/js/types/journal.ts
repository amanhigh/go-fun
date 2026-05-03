export type JournalImage = {
	id: string;
	timeframe: string;
	file_name: string;
	created_at?: string;
};

export type JournalNote = {
	id: string;
	status: string;
	content: string;
	format?: string;
	created_at?: string;
};

export type JournalTag = {
	id: string;
	tag: string;
	type?: string;
	override?: string;
	created_at?: string;
};

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

export type JournalNoteRequest = {
	status: string;
	content: string;
	format: 'MARKDOWN' | 'PLAINTEXT';
};

export type JournalTagRequest = {
	tag: string;
	type: string;
	override?: string;
};

export type JournalFilterKey = 'ticker' | 'type' | 'status' | 'sequence' | 'createdAfter' | 'createdBefore' | 'reviewed' | 'sortBy' | 'sortOrder';

export type JournalFilters = Record<JournalFilterKey, string>;

export type JournalListRequest = Partial<JournalFilters>;

type FilterConfigEntry = {
	queryKey?: string;
	aliases?: readonly string[];
};

const journalFilterConfig: Record<JournalFilterKey, FilterConfigEntry> = {
	ticker: {
		queryKey: 'search',
		aliases: ['ticker'],
	},
	type: {},
	status: {},
	sequence: {},
	createdAfter: {
		queryKey: 'created-after',
	},
	createdBefore: {
		queryKey: 'created-before',
	},
	reviewed: {},
	sortBy: {
		queryKey: 'sort-by',
	},
	sortOrder: {
		queryKey: 'sort-order',
	},
};

export const journalFields: JournalFilterKey[] = ['ticker', 'type', 'status', 'sequence', 'createdAfter', 'createdBefore', 'reviewed', 'sortBy', 'sortOrder'];

export const journalQueryMap: Partial<Record<JournalFilterKey, string>> = journalFields.reduce((queryMap, field) => {
	const entry = journalFilterConfig[field];
	if (!entry.queryKey) return queryMap;
	return {
		...queryMap,
		[field]: entry.queryKey,
	};
}, {} as Partial<Record<JournalFilterKey, string>>);

export const journalReverseMap: Record<string, JournalFilterKey> = journalFields.reduce((reverseMap, field) => {
	const queryKey = journalQueryMap[field] ?? field;
	const aliases = journalFilterConfig[field].aliases ?? [];

	return {
		...reverseMap,
		[queryKey]: field,
		...aliases.reduce<Record<string, JournalFilterKey>>((aliasMap, alias) => ({
			...aliasMap,
			[alias]: field,
		}), {}),
	};
}, {} as Record<string, JournalFilterKey>);

export const journalFilterUrlMapping = {
	fields: journalFields,
	queryMap: journalQueryMap,
	reverseMap: journalReverseMap,
} as const;

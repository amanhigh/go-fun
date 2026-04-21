import type { Journal } from '../client/journal';
import type { ReviewPreset } from './presets';
import type { JournalFilterState } from './filter';
import type { PaginationState } from './pagination';

export type JournalPageState = {
	journals: Journal[];
	reviewPresets: ReviewPreset[];
	activeReviewPreset: string;
	pagination: PaginationState;
	filter: JournalFilterState;
	requestCounter: number;
	loading: boolean;
	errorMessage: string;
};

type CreateJournalPageStateInput = {
	filter: JournalFilterState;
	pagination: PaginationState;
	reviewPresets: ReviewPreset[];
};

export function createJournalPageState(input: CreateJournalPageStateInput): JournalPageState {
	return {
		journals: [],
		reviewPresets: input.reviewPresets,
		activeReviewPreset: '',
		pagination: input.pagination,
		filter: input.filter,
		requestCounter: 0,
		loading: false,
		errorMessage: '',
	};
}

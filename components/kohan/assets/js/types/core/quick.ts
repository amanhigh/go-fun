import type { JournalType, JournalStatus } from '../api/journal/enums';

export interface QuickFilterButton<T extends string = string> {
	readonly label: string;
	readonly className: string;
	toggle(): void;
}

export type QuickConcern = {
	type: QuickFilterButton<JournalType>;
	status: QuickFilterButton<JournalStatus>;
};

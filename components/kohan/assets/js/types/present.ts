import type { JournalTag } from './journal_api';

export type DisplayValue = string | null | undefined;

export type DisplaySpec = {
	icon?: string;
	text: string;
	class: string;
};

export interface Presenter<T = DisplayValue> {
	spec(value: T): DisplaySpec;
	label(value: T): string;
}

export interface TagPresenter extends Presenter<JournalTag> {
}

export interface DatePresenter {
	format(value: DisplayValue): string;
	formatReviewQueueDate(value: DisplayValue): string;
}

export interface PresentationConcern {
	status: Presenter;
	type: Presenter;
	timeframe: Presenter;
	tag: TagPresenter;
	sequence: Presenter;
	date: DatePresenter;
}

export type DisplayValue = string | null | undefined;

export type DisplaySpec = {
	icon?: string;
	text: string;
	class: string;
};

export interface Presenter {
	spec(value: DisplayValue): DisplaySpec;
	label(value: DisplayValue): string;
}

export interface PresentationConcern {
	status: Presenter;
	type: Presenter;
	timeframe: Presenter;
}

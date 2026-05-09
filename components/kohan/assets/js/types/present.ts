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

export abstract class BasePresenter implements Presenter {
	abstract spec(value: DisplayValue): DisplaySpec;

	label(value: DisplayValue): string {
		const s = this.spec(value);
		return s.icon ? `${s.icon} ${s.text}` : s.text;
	}
}

export interface PresentationConcern {
	status: Presenter;
	type: Presenter;
	timeframe: Presenter;
}

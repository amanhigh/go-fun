export type QuickButtonResult = {
	label: string;
	className: string;
	nextValue: string;
};

export interface QuickButton {
	resolve(currentValue: string): QuickButtonResult;
}

export interface QuickFilterButton {
	readonly label: string;
	readonly className: string;
	readonly nextValue: string;
	toggle(): void;
}

export type QuickConcern = {
	type: QuickFilterButton;
	status: QuickFilterButton;
};

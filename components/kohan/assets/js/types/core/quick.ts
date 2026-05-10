export type QuickButtonResult = {
	label: string;
	className: string;
	nextValue: string;
};

export interface QuickButton {
	resolve(currentValue: string): QuickButtonResult;
}

export interface QuickFilterButton {
	button(): QuickButtonResult;
	toggle(): void;
}

export type QuickConcern = {
	status: QuickFilterButton;
};

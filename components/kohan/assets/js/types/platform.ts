export type AlpineGlobal<TFactory extends (...args: never[]) => unknown> = {
	data(name: string, callback: TFactory): void;
};

export type AlpineRefs = {
	reasonTagOverride?: {
		focus?: () => void;
	};
};

export type AlpineContext = {
	$nextTick(callback: () => void): void;
	$refs?: AlpineRefs;
};

export type BrowserWindow = Window & typeof globalThis;

export type AlpineStatic = {
	data<T>(name: string, callback: () => T): void;
};

declare global {
	interface Window {
		Alpine: AlpineStatic;
	}
}

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

export {};

export type AlpineStatic = {
	data<T>(name: string, callback: () => T): void;
};

declare global {
	interface Window {
		Alpine: AlpineStatic;
	}
}



export {};

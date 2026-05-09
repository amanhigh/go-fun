import type { DisplaySpec } from './present';

export type QuickAction = {
	id: string;
	isActive(): boolean;
	display: DisplaySpec;
	apply(): Promise<void>;
};

import type { DisplaySpec } from './presentation_concern';

export type QuickAction = {
	id: string;
	isActive(): boolean;
	display: DisplaySpec;
	apply(): Promise<void>;
};

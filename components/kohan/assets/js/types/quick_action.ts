import type { DisplaySpec } from './presentation_concern';

export type QuickAction = {
	id: string;
	display(): DisplaySpec;
	apply(): Promise<void>;
};

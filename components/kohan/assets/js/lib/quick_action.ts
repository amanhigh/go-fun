import type { DisplaySpec } from '../types/present';
import type { QuickAction } from '../types/quick_action';

export function createQuickAction(
	id: string,
	isActive: () => boolean,
	display: DisplaySpec,
	apply: () => Promise<void>,
): QuickAction {
	return {
		id,
		isActive,
		display,
		apply,
	};
}
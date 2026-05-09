import type { DisplaySpec } from '../types/presentation_concern';
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
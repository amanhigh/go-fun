import type { DisplaySpec } from '../types/presentation_concern';
import type { QuickAction } from '../types/quick_action';

/**
 * Creates a QuickAction with an independent isActive check.
 * display is always a ready spec; isActive decides whether to show it.
 */
export function createQuickAction(
	id: string,
	isActive: () => boolean,
	display: DisplaySpec,
	apply: () => Promise<void>,
): QuickAction {
	return {
		id,
		isActive,
		display: () => display,
		apply,
	};
}

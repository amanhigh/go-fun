import type { DisplaySpec, Presenter } from '../../types/present';
import { BasePresenter } from './base';

const statusCatalog: Record<string, DisplaySpec> = {
	SUCCESS: { icon: '✅', text: 'SUCCESS', class: 'journal-status-success' },
	FAIL: { icon: '❌', text: 'FAIL', class: 'journal-status-fail' },
	RUNNING: { icon: '🏃', text: 'RUNNING', class: 'journal-status-running' },
	SET: { icon: '🎯', text: 'SET', class: 'journal-status-set' },
	JUST_LOSS: { icon: '💔', text: 'JUST_LOSS', class: 'journal-status-just-loss' },
	BROKEN: { icon: '💥', text: 'BROKEN', class: 'journal-status-broken' },
	MISSED: { icon: '🚫', text: 'MISSED', class: 'journal-status-missed' },
};

const fallbackSpec: DisplaySpec = { icon: '🏷️', text: 'Unknown', class: 'journal-display-default' };

class StatusPresenterImpl extends BasePresenter {
	protected catalog = statusCatalog;
	protected fallbackSpec = fallbackSpec;
}

export function NewStatusPresenter(): Presenter {
	return new StatusPresenterImpl();
}

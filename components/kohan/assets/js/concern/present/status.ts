import type { DisplaySpec, DisplayValue, Presenter } from '../../types/present';
import { BasePresenter } from '../../types/present';

const statusCatalog: Record<string, DisplaySpec> = {
	SUCCESS: { icon: '✅', text: 'SUCCESS', class: 'journal-status-success' },
	FAIL: { icon: '❌', text: 'FAIL', class: 'journal-status-fail' },
	RUNNING: { icon: '🏃', text: 'RUNNING', class: 'journal-status-running' },
	SET: { icon: '🎯', text: 'SET', class: 'journal-status-set' },
	JUST_LOSS: { icon: '💔', text: 'JUST_LOSS', class: 'journal-status-just-loss' },
	BROKEN: { icon: '💥', text: 'BROKEN', class: 'journal-status-broken' },
	MISSED: { icon: '🚫', text: 'MISSED', class: 'journal-status-missed' },
	REJECTED: { icon: '🏷️', text: 'REJECTED', class: 'journal-status-rejected' },
};

const fallbackSpec: DisplaySpec = { icon: '🏷️', text: 'Unknown', class: 'journal-display-default' };

class StatusPresenterImpl extends BasePresenter {
	spec(value: DisplayValue): DisplaySpec {
		const key = (value ?? '').trim().toUpperCase();
		return statusCatalog[key] ?? { ...fallbackSpec, text: key || fallbackSpec.text };
	}
}

export function NewStatusPresenter(): Presenter {
	return new StatusPresenterImpl();
}

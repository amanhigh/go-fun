import type { DisplaySpec, DisplayValue, Presenter } from '../../types/core/present';
import { BasePresenter } from './base';

const sequenceDisplayMap: Record<string, DisplaySpec> = {
	MWD: { icon: '🕐', text: 'MWD', class: '' },
	YR: { icon: '📅', text: 'YR', class: '' },
	WDH: { icon: '📅', text: 'WDH', class: '' },
};

class SequencePresenterImpl extends BasePresenter {
	protected catalog = sequenceDisplayMap;
	protected fallbackSpec: DisplaySpec = { icon: '📅', text: '', class: '' };

	spec(value: DisplayValue): DisplaySpec {
		const key = value ?? '';
		if (!key) return { text: '', class: '' };
		return this.catalog[key] ?? { icon: '📅', text: key, class: '' };
	}
}

export function NewSequencePresenter(): Presenter {
	return new SequencePresenterImpl();
}

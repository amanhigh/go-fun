import type { DisplaySpec, DisplayValue, Presenter } from '../../types/present';
import { BasePresenter } from '../../types/present';

const typeCatalog: Record<string, DisplaySpec> = {
	TAKEN: { icon: '📈', text: 'TAKEN', class: 'journal-type-taken' },
	REJECTED: { icon: '📉', text: 'REJECTED', class: 'journal-type-rejected' },
	RESULT: { icon: '🏷️', text: 'RESULT', class: 'journal-type-result' },
	SET: { icon: '🏷️', text: 'SET', class: 'journal-type-set' },
};

const fallbackSpec: DisplaySpec = { icon: '🏷️', text: 'Unknown', class: 'journal-display-default' };

class TypePresenterImpl extends BasePresenter {
	spec(value: DisplayValue): DisplaySpec {
		const key = (value ?? '').trim().toUpperCase();
		return typeCatalog[key] ?? { ...fallbackSpec, text: key || fallbackSpec.text };
	}
}

export function NewTypePresenter(): Presenter {
	return new TypePresenterImpl();
}

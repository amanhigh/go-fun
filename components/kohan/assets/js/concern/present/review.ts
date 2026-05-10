import type { DisplaySpec, DisplayValue, Presenter } from '../../types/present';
import { BasePresenter } from './base';

function formatDate(value: string): string {
	const parsed = new Date(value);
	return Number.isNaN(parsed.getTime()) ? '—' : parsed.toLocaleString();
}

class ReviewPresenterImpl extends BasePresenter {
	protected catalog: Record<string, DisplaySpec> = {};
	protected fallbackSpec: DisplaySpec = { icon: '⏳', text: 'Pending Review', class: '' };

	spec(value: DisplayValue): DisplaySpec {
		if (!value) {
			return { ...this.fallbackSpec };
		}
		return { icon: '✅', text: formatDate(value), class: '' };
	}
}

export function NewReviewPresenter(): Presenter<DisplayValue> {
	return new ReviewPresenterImpl();
}

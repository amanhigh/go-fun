import type { DisplaySpec, DisplayValue, Presenter } from '../../types/present';

function formatDate(value: string): string {
	const parsed = new Date(value);
	return Number.isNaN(parsed.getTime()) ? '—' : parsed.toLocaleString();
}

class ReviewPresenterImpl implements Presenter<DisplayValue> {
	spec(value: DisplayValue): DisplaySpec {
		if (!value) {
			return { icon: '⏳', text: 'Pending Review', class: '' };
		}
		return { icon: '✅', text: formatDate(value), class: '' };
	}

	label(value: DisplayValue): string {
		const s = this.spec(value);
		return s.icon ? `${s.icon} ${s.text}` : s.text;
	}
}

export function NewReviewPresenter(): Presenter<DisplayValue> {
	return new ReviewPresenterImpl();
}

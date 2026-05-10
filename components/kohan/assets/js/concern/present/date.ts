import type { DatePresenter, DisplayValue } from '../../types/present';

const shortMonthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

class DatePresenterImpl implements DatePresenter {
	humanDate(date: Date): string {
		const year = date.getFullYear();
		const month = `${date.getMonth() + 1}`.padStart(2, '0');
		const day = `${date.getDate()}`.padStart(2, '0');
		return `${year}-${month}-${day}`;
	}

	format(value: DisplayValue): string {
		if (!value) return '—';
		const parsed = new Date(value);
		return Number.isNaN(parsed.getTime()) ? '—' : parsed.toLocaleString();
	}

	formatReviewQueueDate(value: DisplayValue): string {
		if (!value) return '—';
		const parsed = new Date(value);
		if (Number.isNaN(parsed.getTime())) return '—';
		const day = parsed.getUTCDate();
		const month = shortMonthNames[parsed.getUTCMonth()] ?? '—';
		const year = `${parsed.getUTCFullYear()}`.slice(-2);
		return `${day} ${month}, '${year}`;
	}
}

export function NewDatePresenter(): DatePresenter {
	return new DatePresenterImpl();
}

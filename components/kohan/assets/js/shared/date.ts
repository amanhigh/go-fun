export function formatTimestamp(value: string | null | undefined): string {
	if (!value) return '—';
	const parsed = new Date(value);
	return Number.isNaN(parsed.getTime()) ? '—' : parsed.toLocaleString();
}

export function formatDateInputValue(date: Date): string {
	const year = date.getFullYear();
	const month = `${date.getMonth() + 1}`.padStart(2, '0');
	const day = `${date.getDate()}`.padStart(2, '0');
	return `${year}-${month}-${day}`;
}

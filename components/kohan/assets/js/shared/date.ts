export function formatTimestamp(value: string | null | undefined): string {
	if (!value) return '—';
	const parsed = new Date(value);
	return Number.isNaN(parsed.getTime()) ? '—' : parsed.toLocaleString();
}

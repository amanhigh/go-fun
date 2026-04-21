export function normalizeTag(value: string): string {
	return (value ?? '').trim().toUpperCase();
}

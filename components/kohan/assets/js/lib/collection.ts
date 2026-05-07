type Identifiable = {
	id: string;
};

export function prependById<T extends Identifiable>(items: T[], item: T): T[] {
	return [item, ...items.filter((candidate) => candidate.id !== item.id)];
}

export function removeById<T extends Identifiable>(items: T[], itemId: string): T[] {
	return items.filter((item) => item.id !== itemId);
}

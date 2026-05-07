import { normalizeTag } from './tags';

export type ManagementTagPreset = {
	value: string;
	label: string;
	tone: string;
};

export const managementTagPresets: readonly ManagementTagPreset[] = [
	{ value: 'ntr', label: 'NTR', tone: 'emerald' },
	{ value: 'enl', label: 'ENL', tone: 'sky' },
	{ value: 'slt', label: 'SLT', tone: 'rose' },
	{ value: 'fz', label: 'FZ', tone: 'violet' },
	{ value: 'nbe', label: 'NBE', tone: 'amber' },
	{ value: 'ws', label: 'WS', tone: 'slate' },
	{ value: 'important', label: 'IMPORTANT', tone: 'fuchsia' },
	{ value: 'be', label: 'BE', tone: 'orange' },
];

const toneMap: Record<string, string> = {};
for (const preset of managementTagPresets) {
	toneMap[normalizeTag(preset.value)] = preset.tone;
}

export function managementTagTone(value: string): string {
	return toneMap[normalizeTag(value)] ?? 'slate';
}

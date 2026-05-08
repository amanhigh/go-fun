// ===== Types =====

export type ManagementTagPreset = {
	value: string;
	label: string;
	tone: string;
};

// ===== Presets =====

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

// ===== Tone Lookup =====

const DEFAULT_MANAGEMENT_TONE = 'slate';

export function managementTagTone(value: string): string {
	const preset = managementTagPresets.find((p) => p.value === value);
	return preset?.tone ?? DEFAULT_MANAGEMENT_TONE;
}

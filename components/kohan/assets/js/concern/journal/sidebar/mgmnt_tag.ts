import { createAsyncFeedback } from '../../../lib/async_feedback';
import { normalizeTag } from '../../../lib/tags';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

type ManagementTagPreset = {
	value: string;
	label: string;
	tone: string;
};

const managementTagPresets: ManagementTagPreset[] = [
	{ value: 'ntr', label: 'NTR', tone: 'emerald' },
	{ value: 'enl', label: 'ENL', tone: 'sky' },
	{ value: 'slt', label: 'SLT', tone: 'rose' },
	{ value: 'fz', label: 'FZ', tone: 'violet' },
	{ value: 'nbe', label: 'NBE', tone: 'amber' },
	{ value: 'ws', label: 'WS', tone: 'slate' },
	{ value: 'important', label: 'IMPORTANT', tone: 'fuchsia' },
	{ value: 'be', label: 'BE', tone: 'orange' },
];

export function MgmntTagConcern(pg: JournalDetailPageProvider) {
	return {
		...createAsyncFeedback(),
		presets: managementTagPresets,
		pendingValue: '',

		showTakenTags() {
			return normalizeTag(pg().current.journal?.type ?? '') === 'TAKEN';
		},
		hasTag(value: string) {
			const normalizedValue = normalizeTag(value);
			return pg().sidebar.tags.management().some((tag: JournalTag) => normalizeTag(tag.tag ?? '') === normalizedValue);
		},
		buttonClass(preset: ManagementTagPreset) {
			const tagKey = normalizeTag(preset.value);
			const isActive = this.hasTag(preset.value);
			const isPending = this.submitting && normalizeTag(this.pendingValue) === tagKey;
			const baseClass = isActive ? `journal-management-active-${preset.tone}` : `journal-management-base-${preset.tone}`;
			return isPending ? `journal-management-pending ${baseClass}` : baseClass;
		},
		async submit(tagValue: string) {
			if (!pg().current.journal || this.submitting) return;
			this.pendingValue = tagValue;
			await this.run(
				() => this.addTag(tagValue),
				`${normalizeTag(tagValue)} tag added.`,
				'Unable to save management tag.',
			);
			this.pendingValue = '';
		},

		async addTag(tagValue: string) {
			const payload: JournalTagRequest = { tag: tagValue, type: 'MANAGEMENT' };
			const page = pg();
			const envelope = await page.tagClient.create(page.current.journalId, payload);
			page.sidebar.tags.prepend(envelope.data as JournalTag);
		},
	};
}

import type { DisplaySpec, TagPresenter } from '../../types/present';
import type { JournalTag } from '../../types/journal_api';
import { normalizeTag } from '../../lib/tags';

const reasonIcons = { trend: '📈', default: '⚡' };

const managementDisplayMap: Record<string, DisplaySpec> = {
	ntr: { icon: '🔧', text: 'NTR', class: 'journal-management-base-emerald' },
	enl: { icon: '🔧', text: 'ENL', class: 'journal-management-base-sky' },
	slt: { icon: '🔧', text: 'SLT', class: 'journal-management-base-rose' },
	fz: { icon: '🔧', text: 'FZ', class: 'journal-management-base-violet' },
	nbe: { icon: '🔧', text: 'NBE', class: 'journal-management-base-amber' },
	ws: { icon: '🔧', text: 'WS', class: 'journal-management-base-slate' },
	important: { icon: '⭐', text: 'Important', class: 'journal-management-base-fuchsia' },
	be: { icon: '🔧', text: 'BE', class: 'journal-management-base-orange' },
};

function reasonTagText(tag: JournalTag): string {
	const name = tag.tag ?? '';
	const override = tag.override ? ` → ${tag.override}` : '';
	return `${name}${override}`;
}

function fallbackTagText(tag: JournalTag): string {
	return [tag.tag, tag.override, tag.type].filter(Boolean).join(' • ');
}

class TagPresenterImpl implements TagPresenter {
	spec(tag: JournalTag): DisplaySpec {
		const type = normalizeTag(tag.type ?? '');
		switch (type) {
			case 'REASON':
				return this.reasonSpec(tag);
			case 'DIRECTION':
				return this.directionSpec(tag);
			case 'MANAGEMENT':
				return this.managementTagSpec(tag);
			default:
				return this.fallbackSpec(tag);
		}
	}

	label(tag: JournalTag): string {
		const s = this.spec(tag);
		return s.icon ? `${s.icon} ${s.text}` : s.text;
	}

	private reasonSpec(tag: JournalTag): DisplaySpec {
		const name = tag.tag ?? '';
		const icon = name.toLowerCase().includes('trend') ? reasonIcons.trend : reasonIcons.default;
		return { icon, text: reasonTagText(tag), class: '' };
	}

	private directionSpec(tag: JournalTag): DisplaySpec {
		const name = tag.tag ?? '';
		return { icon: '🔀', text: name, class: '' };
	}

	private managementTagSpec(tag: JournalTag): DisplaySpec {
		const key = normalizeTag(tag.tag ?? '');
		return managementDisplayMap[key] ?? { icon: '🔧', text: key, class: '' };
	}

	private fallbackSpec(tag: JournalTag): DisplaySpec {
		return { icon: '🏷️', text: fallbackTagText(tag), class: '' };
	}
}

export function NewTagPresenter(): TagPresenter {
	return new TagPresenterImpl();
}

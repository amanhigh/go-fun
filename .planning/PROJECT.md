# Go Fun / Kohan

## What This Is

Kohan is a journal workflow inside the Go Fun project. This milestone focuses on the journal detail page so users can review an entry, add tags, and add notes without leaving the page.

## Core Value

Make journal review fast and reliable from the detail page.

## Requirements

### Validated

- `JournalDetailPage` can display a single journal entry with images, tags, and notes.
- Existing journal APIs already support fetching a journal and updating review status.

### Active

- [ ] Mark a journal as reviewed today from the detail page.
- [ ] Add one or more tags with a quick-add interaction.
- [ ] Add a note with an inline composer.
- [ ] Reload the page after each save so the latest data is visible.

### Out of Scope

- Bulk edit of multiple journals — this milestone is single-entry only.
- Separate tag management pages — tags are edited in-place only.
- Rich text notes — notes remain plain text.
- Search/filtering on the detail page — not needed for this milestone.

## Context

- Existing detail UI lives in `components/kohan/ui/pages/journal_detail.templ`.
- Client-side behavior is already Alpine-based in `components/kohan/assets/js/app.js`.
- Journal review update APIs already exist in `components/kohan/handler/journal.go`.
- The codebase uses handler → manager → repository flow and Templ + Tailwind UI patterns.

## Constraints

- **Architecture**: Keep handler → manager → repository separation intact.
- **UI**: Reuse existing TemplUI/Tailwind patterns and avoid unnecessary wrapper elements.
- **Behavior**: Save actions should reload the page after completion.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Review is one-click “Mark reviewed today” | Fastest path for the milestone | ✓ Good |
| Tags are inline quick-add and can be added multiple at once | Matches the existing detail-page workflow | ✓ Good |
| Notes use an inline quick composer | Keeps editing close to the existing notes section | ✓ Good |
| Save actions reload the page | Simplifies state sync after edits | ✓ Good |

---
*Last updated: 2026-04-16 after Milestone v1.0 start*

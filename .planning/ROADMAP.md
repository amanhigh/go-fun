# Roadmap: Go Fun / Kohan

**Milestone:** v1.0 — Journal Update/Edit Page

## Phase 1 — Review Action

**Goal:** Let the user mark a journal as reviewed today from the detail page.

**Requirements:** `JDET-01`

**Success criteria:**
1. The detail page shows a review action for an unreviewed journal.
2. Clicking the action sends the review update and saves today’s date.
3. The reviewed state is visible after the page refreshes.

## Phase 2 — Tag Quick Add

**Goal:** Let the user add one or more tags inline on the journal detail page.

**Requirements:** `JDET-02`

**Success criteria:**
1. The detail page exposes a quick-add tag control near existing tags.
2. Multiple tags can be added without leaving the page.
3. Saved tags are visible after refresh.

## Phase 3 — Note Composer

**Goal:** Let the user add a note inline on the journal detail page.

**Requirements:** `JDET-03`

**Success criteria:**
1. The notes section includes an inline composer.
2. Submitting a note saves it against the journal.
3. The new note appears after refresh.

## Phase 4 — Save/Refresh Behavior

**Goal:** Keep the detail page in sync by reloading after edits.

**Requirements:** `JDET-04`

**Success criteria:**
1. Every save action ends with a page reload.
2. Reloaded state reflects the latest review, tags, and notes.
3. Edit flows remain consistent with the existing Templ + Alpine page structure.

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| JDET-01 | Phase 1 | Pending |
| JDET-02 | Phase 2 | Pending |
| JDET-03 | Phase 3 | Pending |
| JDET-04 | Phase 4 | Pending |

## Notes

- Phase numbering starts at 1 for this new milestone.
- The roadmap intentionally keeps each edit flow isolated so planning and implementation can proceed incrementally.

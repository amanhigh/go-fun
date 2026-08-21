-- 000003 down: Convert journals.top_timeframe back to journals.sequence AND
-- drop images.image_type.
--
-- Part A (journal top_timeframe):
-- Reverse mapping (new -> legacy). This mapping is LOSSY because the legacy
-- domain ('MWD','YR','WDH') cannot represent the new domain
-- ('YR','SMN','TMN','MN') distinctly:
--   YR  -> YR   (unchanged)
--   SMN -> YR   (SMN and YR both collapse to legacy YR)
--   TMN -> MWD
--   MN  -> WDH
-- Rows that were originally YR become SMN on the way up and YR on the way
-- down; rows that were originally MWD/WDH round-trip exactly. This is the
-- best-effort reverse required by rollback semantics.
--
-- Simple in-place inverse of the up migration: add sequence, backfill it from
-- top_timeframe, then drop top_timeframe. No table rebuild is needed.
--
-- Part B (image classification):
-- Drop images.image_type before the journal rollback so the schema matches the
-- pre-migration state.

-- 0. Drop images.image_type first so the schema is restored to its pre-000003
--    state before the journal rollback runs.
ALTER TABLE images DROP COLUMN image_type;

-- 1. Add sequence. The DEFAULT 'YR' is required because SQLite only allows a
--    NOT NULL column to be added with a non-NULL default; the UPDATE in step 2
--    immediately overwrites it with the reverse mapping.
ALTER TABLE journals ADD COLUMN sequence VARCHAR(3) NOT NULL DEFAULT 'YR' CHECK (sequence IN ('MWD','YR','WDH'));

-- 2. Backfill from top_timeframe using the reverse mapping.
UPDATE journals SET sequence = CASE top_timeframe
    WHEN 'YR' THEN 'YR'
    WHEN 'SMN' THEN 'YR'
    WHEN 'TMN' THEN 'MWD'
    WHEN 'MN' THEN 'WDH'
END;

-- 3. Drop the new column. Steps 1-3 run atomically inside the migration
--    transaction, so the table is never observed in an intermediate state.
ALTER TABLE journals DROP COLUMN top_timeframe;

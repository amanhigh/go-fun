-- 000003: Convert journals.sequence to journals.top_timeframe
--
-- Authoritative mapping (legacy -> new):
--   YR  -> SMN
--   MWD -> TMN
--   WDH -> MN
--
-- Simple in-place SQLite migration: add top_timeframe, backfill it from
-- sequence, then drop sequence. No table rebuild is needed because the child
-- tables (images/tags/notes) are untouched.

-- 1. Add top_timeframe. The DEFAULT 'SMN' is required because SQLite only
--    allows NOT NULL columns to be added with a non-NULL default; the UPDATE
--    in step 2 immediately overwrites it with the authoritative mapping.
ALTER TABLE journals ADD COLUMN top_timeframe VARCHAR(3) NOT NULL DEFAULT 'SMN' CHECK (top_timeframe IN ('YR', 'SMN', 'TMN', 'MN'));

-- 2. Backfill from sequence using the authoritative mapping.
UPDATE journals SET top_timeframe = CASE sequence
    WHEN 'YR' THEN 'SMN'
    WHEN 'MWD' THEN 'TMN'
    WHEN 'WDH' THEN 'MN'
END;

-- 3. Drop the legacy column. Steps 1-3 run atomically inside the migration
--    transaction, so the table is never observed in an intermediate state.
ALTER TABLE journals DROP COLUMN sequence;
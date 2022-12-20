-- Add new columns without constraints first, migrate, then run
-- Go code to recalculate fingerprints and perform the followup
-- migration which will add constraints.

ALTER TABLE pubkeys ADD COLUMN
  type TEXT NOT NULL DEFAULT '';

ALTER TABLE pubkeys ADD COLUMN
  fingerprint_legacy TEXT NOT NULL DEFAULT '';

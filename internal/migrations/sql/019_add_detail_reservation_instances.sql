ALTER TABLE reservation_instances ADD COLUMN
  detail JSONB NOT NULL DEFAULT '{}'::jsonb

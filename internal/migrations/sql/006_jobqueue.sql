-- We use pg_catalog because we have two schemas: public and integration
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA "pg_catalog";

-- https://github.com/osbuild/osbuild-composer/blob/main/internal/jobqueue/dbjobqueue/schemas/001_jobs.sql
CREATE TABLE jobs(
        id uuid PRIMARY KEY,
        token uuid,
        type varchar NOT NULL,
        args jsonb,
        result jsonb,
        queued_at timestamp NOT NULL,
        started_at timestamp,
        finished_at timestamp,
        canceled boolean NOT NULL DEFAULT FALSE,

        CONSTRAINT not_finished_when_not_started
          CHECK (finished_at IS NULL OR started_at IS NOT NULL),

        CONSTRAINT not_finished_when_canceled
          CHECK (CANCELED = FALSE OR finished_at IS NULL),

        CONSTRAINT chronologic_started_at
          CHECK (started_at IS NULL OR queued_at <= started_at),

        CONSTRAINT chronologic_finished_at
          CHECK (finished_at IS NULL OR started_at <= finished_at),

        CONSTRAINT token_is_never_uuid_nil
          CHECK (token IS NULL OR token != uuid_nil()),

        CONSTRAINT token_is_set_when_started
          CHECK (started_at IS NULL OR token IS NOT NULL)
);

CREATE TABLE job_dependencies(
        job_id uuid REFERENCES jobs(id),
        dependency_id uuid REFERENCES jobs(id)
);

CREATE TABLE heartbeats(
       token uuid PRIMARY KEY,
       id uuid REFERENCES jobs(id),
       heartbeat timestamp NOT NULL,

       CONSTRAINT token_is_never_uuid_nil
         CHECK (token != uuid_nil())
);

CREATE VIEW ready_jobs AS
  SELECT *
  FROM jobs
  WHERE started_at IS NULL
    AND canceled = FALSE
    AND id NOT IN (
      SELECT job_id
      FROM job_dependencies JOIN jobs ON dependency_id = id
      WHERE finished_at IS NULL
    )
  ORDER BY queued_at ASC;

-- https://github.com/osbuild/osbuild-composer/blob/main/internal/jobqueue/dbjobqueue/schemas/002_channels.sql

ALTER TABLE jobs
  ADD COLUMN channel varchar NOT NULL DEFAULT '';

CREATE OR REPLACE VIEW ready_jobs AS
SELECT *
FROM jobs
WHERE started_at IS NULL
  AND canceled = FALSE
  AND id NOT IN (
  SELECT job_id
  FROM job_dependencies JOIN jobs ON dependency_id = id
  WHERE finished_at IS NULL
)
ORDER BY queued_at ASC;

-- https://github.com/osbuild/osbuild-composer/blob/main/internal/jobqueue/dbjobqueue/schemas/003_dependencies_cascade.sql

ALTER TABLE job_dependencies

  DROP CONSTRAINT job_dependencies_dependency_id_fkey,

  DROP CONSTRAINT job_dependencies_job_id_fkey,

  ADD CONSTRAINT job_dependencies_dependency_id_fkey
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,

  ADD CONSTRAINT job_dependencies_job_id_fkey
    FOREIGN KEY (dependency_id) REFERENCES jobs(id) ON DELETE CASCADE;

-- https://github.com/osbuild/osbuild-composer/blob/main/internal/jobqueue/dbjobqueue/schemas/004_heartbeats_cascade.sql

ALTER TABLE heartbeats

  DROP CONSTRAINT heartbeats_id_fkey,

  ADD CONSTRAINT heartbeats_id_fkey
    FOREIGN KEY (id) REFERENCES jobs(id) ON DELETE CASCADE;

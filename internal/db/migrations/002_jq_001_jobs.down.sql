BEGIN;

DROP VIEW ready_jobs;
DROP TABLE heartbeats;
DROP TABLE job_dependencies;
DROP TABLE jobs;

COMMIT;
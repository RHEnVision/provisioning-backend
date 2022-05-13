--
-- Testing seed data, to execute this file during migration set DB_SEED variable to "dev_small".
-- Keep this file idempotent so we can run it over and over again to get
--
BEGIN;

INSERT INTO accounts(account_number, org_id)
VALUES
    ('13', '000013'),
    ('15', '000015'),
    (NULL, '000042'),
    (NULL, '000077')
ON CONFLICT DO NOTHING;

COMMIT;

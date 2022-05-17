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

INSERT INTO pubkeys(account_id, name, body)
VALUES
  ((SELECT id FROM accounts WHERE account_number = '13'), 'lzap-ed25519-2021', 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap')
ON CONFLICT DO NOTHING;

COMMIT;

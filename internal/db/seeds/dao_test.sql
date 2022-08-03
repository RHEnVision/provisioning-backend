--
-- Testing seed data, to execute this file during migration set DB_SEED variable to "dao_test".
--
BEGIN;

INSERT INTO accounts(id, account_number, org_id)
VALUES (1, '1', '1')
  ON CONFLICT DO NOTHING;

INSERT INTO pubkeys(id, account_id, name, body)
VALUES (1, 1, 'lzap-ed25519-2021',
        'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap')
ON CONFLICT DO NOTHING;

-- Reset all primary key sequences. This can possibly slow down seeds in tests, in that case
-- let's use implicit primary keys.
SELECT reset_sequences('integration');

COMMIT;

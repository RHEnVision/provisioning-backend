--
-- Testing seed data, to execute this file during migration set DB_SEED variable to "dev_small".
-- Keep this file idempotent. Always specify primary keys in the range of 1-100.
--
BEGIN;

INSERT INTO accounts(id, account_number, org_id)
VALUES (1, '13', '000013'),
       (2, '15', '000015'),
       (3, NULL, '000042'),
       (4, NULL, '000077')
ON CONFLICT DO NOTHING;

INSERT INTO pubkeys(id, account_id, name, body)
VALUES (1, 1, 'lzap-ed25519-2021',
        'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap')
ON CONFLICT DO NOTHING;

-- Reset all primary key sequences (columns named "id") to the maximum value.
SELECT reset_sequences();

COMMIT;

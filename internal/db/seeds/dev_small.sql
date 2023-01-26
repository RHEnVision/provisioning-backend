--
-- Testing seed data, to execute this file during migration set DB_SEED variable to "dev_small".
-- Keep this file idempotent. Always specify primary keys in the range of 1-100.
--
BEGIN;

-- Seed some account numbers, artificial or stage environment
INSERT INTO accounts(id, account_number, org_id)
VALUES (1, '13', '000013'),       -- non-existing account
       (2, NULL, '000042'),       -- non-existing account
       (3, '6395343', '13446659') -- stage account
ON CONFLICT DO NOTHING;

-- Seed some pubkeys, feel free to add your own key and associate it with your account
INSERT INTO pubkeys(id, account_id, name, body, type, fingerprint, fingerprint_legacy)
VALUES (1, 3, 'lzap-ed25519-2021',
        'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap',
        'ssh-ed25519',
        'gL/y6MvNmJ8jDXtsL/oMmK8jUuIefN39BBuvYw/Rndk=',
        '89:c5:99:b5:33:48:1c:84:be:da:cb:97:45:b0:4a:ee')
ON CONFLICT DO NOTHING;

-- Reset all primary key sequences (columns named "id") to the maximum value.
SELECT reset_sequences('public');

COMMIT;

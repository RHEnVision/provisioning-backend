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
        'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap+edkey@redhat.com',
        'ssh-ed25519',
        'gL/y6MvNmJ8jDXtsL/oMmK8jUuIefN39BBuvYw/Rndk=',
        'ee:f1:d4:62:99:ab:17:d9:3b:00:66:62:32:b2:55:9e'),
(2, 3, 'lzap-rsa-2010',
        'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8w6DONv1qn3IdgxSpkYOClq7oe7davWFqKVHPbLoS6+dFInru7gdEO5byhTih6+PwRhHv/b1I+Mtt5MDZ8Sv7XFYpX/3P/u5zQiy1PkMSFSz0brRRUfEQxhXLW97FJa7l+bej2HJDt7f9Gvcj+d/fNWC9Z58/GX11kWk4SIXaKotkN+kWn54xGGS7Zvtm86fP59Srt6wlklSsG8mZBF7jVUjyhAgm/V5gDFb2/6jfiwSb2HyJ9/NbhLkWNdwrvpdGZqQlYhnwTfEZdpwizW/Mj3MxP5O31HN45aE0wog0UeWY4gvTl4Ogb6kescizAM6pCff3RBslbFxLdOO7cR17 lzap+rsakey@redhat.com',
        'ssh-rsa',
        'ENShRe/0uDLSw9c+7tc9PxkD/p4blyB/DTgBSIyTAJY=',
        '89:c5:99:b5:33:48:1c:84:be:da:cb:97:45:b0:4a:ee')
ON CONFLICT DO NOTHING;

-- Reset all primary key sequences (columns named "id") to the maximum value.
SELECT reset_sequences('public');

COMMIT;

ALTER TABLE aws_reservation_details
DROP CONSTRAINT aws_reservation_details_pubkey_id_fkey,
ADD CONSTRAINT aws_reservation_details_pubkey_id_fkey
FOREIGN key (pubkey_id) REFERENCES pubkeys(id) ON DELETE CASCADE;

ALTER TABLE gcp_reservation_details
DROP CONSTRAINT gcp_reservation_details_pubkey_id_fkey,
ADD CONSTRAINT gcp_reservation_details_pubkey_id_fkey
FOREIGN key (pubkey_id) REFERENCES pubkeys(id) ON DELETE CASCADE;

ALTER TABLE pubkey_resources
DROP CONSTRAINT pubkey_resources_pubkey_id_fkey,
ADD CONSTRAINT pubkey_resources_pubkey_id_fkey
FOREIGN key (pubkey_id) REFERENCES pubkeys(id) ON DELETE CASCADE;

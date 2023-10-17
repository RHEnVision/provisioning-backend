ALTER TABLE aws_reservation_details
DROP CONSTRAINT aws_reservation_details_pubkey_id_fkey,
ADD CONSTRAINT aws_reservation_details_pubkey_id_fkey
FOREIGN key (pubkey_id) REFERENCES pubkeys(id) ON DELETE SET NULL;

ALTER TABLE gcp_reservation_details
DROP CONSTRAINT gcp_reservation_details_pubkey_id_fkey,
ADD CONSTRAINT gcp_reservation_details_pubkey_id_fkey
FOREIGN key (pubkey_id) REFERENCES pubkeys(id) ON DELETE SET NULL;

ALTER TABLE azure_reservation_details
DROP CONSTRAINT azure_reservation_details_pubkey_id_fkey,
ADD CONSTRAINT azure_reservation_details_pubkey_id_fkey
FOREIGN key (pubkey_id) REFERENCES pubkeys(id) ON DELETE SET NULL;

ALTER TABLE aws_reservation_details ALTER COLUMN pubkey_id DROP NOT NULL;
ALTER TABLE gcp_reservation_details ALTER COLUMN pubkey_id DROP NOT NULL;
ALTER TABLE azure_reservation_details ALTER COLUMN pubkey_id DROP NOT NULL;

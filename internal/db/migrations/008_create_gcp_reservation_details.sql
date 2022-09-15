
CREATE TABLE gcp_reservation_details
(
  reservation_id BIGINT NOT NULL,
  provider INTEGER NOT NULL DEFAULT provider_type_gcp(),
  pubkey_id BIGINT NOT NULL REFERENCES pubkeys(id),
  source_id TEXT NOT NULL,
  image_id TEXT NOT NULL,
  gcp_operation_name TEXT,
  detail JSONB,

  FOREIGN KEY (reservation_id, provider) REFERENCES reservations(id, provider) ON DELETE CASCADE,
  CHECK (provider = provider_type_gcp())
);


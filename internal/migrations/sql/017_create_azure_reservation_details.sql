
CREATE TABLE azure_reservation_details
(
  reservation_id BIGINT NOT NULL,
  provider_fk INTEGER NOT NULL DEFAULT provider_type_azure(),
  pubkey_id BIGINT NOT NULL REFERENCES pubkeys(id),
  source_id TEXT NOT NULL,
  image_id TEXT NOT NULL,
  detail JSONB,

  FOREIGN KEY (reservation_id, provider_fk) REFERENCES reservations(id, provider) ON DELETE CASCADE,
  CHECK (provider_fk = provider_type_azure())
);


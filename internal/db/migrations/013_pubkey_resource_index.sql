-- Write your migrate up statements here
DROP INDEX pubkey_resources_pubkey_id_provider;
CREATE UNIQUE INDEX pubkey_resources_pubkey_id_provider ON pubkey_resources(pubkey_id, source_id, region);

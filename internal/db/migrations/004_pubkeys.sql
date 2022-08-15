CREATE TABLE pubkeys
(
  id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  account_id BIGINT NOT NULL REFERENCES accounts(id),
  name TEXT NOT NULL CHECK (NOT empty(name)),
  body TEXT NOT NULL CHECK (NOT empty(body))
);

CREATE TABLE pubkey_resources
(
  id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  tag TEXT NOT NULL DEFAULT '',
  pubkey_id BIGINT NOT NULL REFERENCES pubkeys(id),
  provider INTEGER NOT NULL CHECK (valid_provider(provider)),
  source_id INTEGER NOT NULL,
  handle TEXT NOT NULL
);

CREATE UNIQUE INDEX pubkey_resources_pubkey_id_provider ON pubkey_resources(pubkey_id, provider);

CREATE TRIGGER prevent_tag_update_on_pubkey_resources
  BEFORE UPDATE
  ON pubkey_resources
  FOR EACH ROW
EXECUTE PROCEDURE prevent_tag_update();

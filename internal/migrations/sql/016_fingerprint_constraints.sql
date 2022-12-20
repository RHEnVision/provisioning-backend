ALTER TABLE pubkeys ADD CONSTRAINT
  pubkeys_type_check CHECK (NOT empty(type));

ALTER TABLE pubkeys ADD CONSTRAINT
  pubkeys_fingerprint_legacy_check CHECK (NOT empty(fingerprint_legacy));

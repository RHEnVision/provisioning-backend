-- Actual reservation limit per account and provider type (returns constant number)
CREATE OR REPLACE FUNCTION reservations_rate_limit() RETURNS INTEGER AS
$reservations_rate_limit$
BEGIN
  RETURN 2;
END;
$reservations_rate_limit$ LANGUAGE plpgsql;

-- Rate limiting function (throws exception when exceeded)
CREATE OR REPLACE FUNCTION reservations_rate() RETURNS TRIGGER AS
$reservations_rate$
DECLARE
  maximum INTEGER := reservations_rate_limit();
  last_rec RECORD;
BEGIN
  FOR last_rec IN SELECT COUNT(*) FROM reservations WHERE account_id = NEW.account_id AND provider = NEW.provider AND success IS NULL AND created_at >= now() - INTERVAL '1 second'
    LOOP
      IF last_rec.count >= maximum THEN
        -- When changing the exception string, also change the ErrReservationRateExceeded handling Go code
        RAISE EXCEPTION 'too many pending reservations for this provider (maximum % per second)', maximum;
      END IF;
    END LOOP;

  RETURN NEW;
END;
$reservations_rate$ LANGUAGE plpgsql;

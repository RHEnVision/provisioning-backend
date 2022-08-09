package payloads

import "database/sql"

func SqlNullToStringPtr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	str := s.String
	return &str
}

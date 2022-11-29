package payloads

import "database/sql"

func SqlNullToStringPtr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	str := s.String
	return &str
}

func StringNullToEmpty(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

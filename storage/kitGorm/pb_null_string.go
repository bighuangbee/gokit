package kitGorm

import (
	"database/sql"
	"encoding/json"
)

// PBNullString PBNullString
type PBNullString struct {
	sql.NullString
}

// NewPBNullString new a not null PBNullString from string
func NewPBNullString(s string) PBNullString {
	return PBNullString{
		NullString: sql.NullString{
			String: s,
			Valid:  true,
		},
	}
}

// NullableString 空字符串会处理成NULL
func NullableString(s string) (r PBNullString) {
	r.NullString.String = s
	if s == "" {
		r.Valid = false
	} else {
		r.Valid = true
	}
	return
}

// MarshalJSON MarshalJSON
func (ns PBNullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON UnmarshalJSON
func (ns *PBNullString) UnmarshalJSON(data []byte) error {
	var value string
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	if value != "" {
		ns.NullString.String = value
		ns.Valid = true
	} else {
		ns.NullString.String = ""
		ns.Valid = false
	}
	return nil
}

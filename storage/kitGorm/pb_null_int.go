package kitGorm

import (
	"database/sql"
	"encoding/json"
)

// PBNullString PBNullString
type PBNullInt32 struct {
	sql.NullInt32
}

// NewPBNullString new a not null PBNullString from string
func NewPBPBNullInt32(i int32) PBNullInt32 {
	return PBNullInt32{
		NullInt32: sql.NullInt32{
			Int32: i,
			Valid: true,
		},
	}
}

// MarshalJSON MarshalJSON
func (ns PBNullInt32) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Int32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON UnmarshalJSON
func (ns *PBNullInt32) UnmarshalJSON(data []byte) error {
	var value int32
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	if value != 0 {
		ns.Int32 = value
		ns.Valid = true
	} else {
		ns.Int32 = 0
		ns.Valid = false
	}
	return nil
}

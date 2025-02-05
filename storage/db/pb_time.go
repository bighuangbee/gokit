package db

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

// PBTime PBTime
// 支持json编码解码ptypes/timestamp
type PBTime time.Time

// MarshalJSON Marsha
func (t PBTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t))
}

// UnmarshalJSON unmarshal
func (t *PBTime) UnmarshalJSON(data []byte) error {
	var p timestamp.Timestamp
	if err := json.Unmarshal(data, &p); err == nil {
		*t = PBTime(p.AsTime())
		return nil
	}
	return json.Unmarshal(data, (*time.Time)(t))
}

// String for sql log, print readable format
func (t PBTime) String() string {
	ts := time.Time(t)
	return ts.Format(time.RFC3339)
}

func (t PBTime) Unix() int64 {
	return time.Time(t).Unix()
}

// Value return json value, implement driver.Valuer interface
// Value insert into database conversion
func (t PBTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan scan value into Jsonb, implements sql.Scanner interface
// Scan read from database conversion
func (t *PBTime) Scan(src interface{}) error {
	if val, ok := src.(time.Time); ok {
		*t = PBTime(val)
	}
	return nil
}

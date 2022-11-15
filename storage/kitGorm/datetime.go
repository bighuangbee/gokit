package kitGorm

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time time.Time

func init() {
	sh, _ := time.LoadLocation(Zone)
	time.Local = sh
}

var (
	TimeFormart = "2006-01-02 15:04:05"
	Zone        = "Asia/Shanghai"
	DateZero    = "2006-01-02 00:00:00"
)

// UnmarshalJSON implements json unmarshal interface.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormart+`"`, string(data), time.Local)
	fmt.Println("------now",now, Time(now))
	*t = Time(now)
	return
}

// MarshalJSON implements json marshal interface.
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormart)+2)
	b = append(b, '"')
	if t != Time(time.Time{}){ //返回空字符串
		b = time.Time(t).AppendFormat(b, TimeFormart)
	}
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	ti := Time{}
	if t == ti {
		return ""
	}
	return t.Local().Format(TimeFormart)
}

func (t Time) Local() time.Time {
	loc, _ := time.LoadLocation(Zone)
	return time.Time(t).In(loc)
}

func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil	//Time的默认值, Writor
	}
	return ti, nil
}

func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func TimeZero() string{
	return "0000-00-00 00:00:00"
}

func TimeNil() Time {
	return Time(time.Time{})
}

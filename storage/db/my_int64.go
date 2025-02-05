package db

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//	MyInt64 自定义
//
// 支持前端js string和后端c++、java的int64
type MyInt64 int64

func (t *MyInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的decimal字符串
	str := string(data)
	str = strings.Trim(str, "\"")
	d, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil
	}
	*t = MyInt64(d)
	return err
}

func (t MyInt64) MarshalJSON() ([]byte, error) {
	str := fmt.Sprint(int64(t))
	return []byte(str), nil
}

func (t MyInt64) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t *MyInt64) Scan(v interface{}) error {
	switch vt := v.(type) {
	case int64:
		*t = MyInt64(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *MyInt64) String() string {
	return fmt.Sprint(int64(*t))
}
func (ts *MyInt64) ToInt64() int64 {
	return int64(*ts)
}

package kitGorm

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//  MyUint64 自定义
// 支持前端js string和后端c++、java的int64
type MyUint64 uint64

func (t *MyUint64) UnmarshalJSON(data []byte) error {
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
	*t = MyUint64(d)
	return err
}

func (t MyUint64) MarshalJSON() ([]byte, error) {
	str := fmt.Sprint(uint64(t))
	return []byte(str), nil
}

func (t MyUint64) Value() (driver.Value, error) {
	return uint64(t), nil
}

func (t *MyUint64) Scan(v interface{}) error {
	/* typeOfResp := reflect.TypeOf(v)
	fmt.Println(fmt.Sprintf(" resp type is %s, kind is %s",
		typeOfResp, typeOfResp.Kind())) */
	switch vt := v.(type) {
	case uint64:
		*t = MyUint64(vt)
	case int64:
		*t = MyUint64(vt)
	case []uint8: //db是 bigint unsigned 18446744073709551615
		// 字符串转成 time.Time 类型
		str := string(vt)
		d, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*t = MyUint64(d)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *MyUint64) String() string {
	return fmt.Sprint(uint64(*t))
}
func (ts *MyUint64) ToUint64() uint64 {
	return uint64(*ts)
}

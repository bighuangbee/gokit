package tools

import (
	"fmt"
	"testing"
	"time"
)


func TestConvertString(t *testing.T) {
	str := ConvertString("aaa达到", "GB2312")
	fmt.Println(str)
}

type UpdateTimeEmbed struct {
	UpdatedAt time.Time `gorm:"column:update_time" json:"update_time" form:"update_time"`
}

type IdEmbed struct {
	ID int64 `gorm:"column:id" json:"id" form:"id"  dbupdate:"-"`
}

type User struct {
	IdEmbed
	Name string `json:"name" dbupdate:"name"`
	TaskId int64 `json:"task_id" dbupdate:"task_id,omitempty"`
	Nickname string `json:"nickname" dbupdate:"nickname"`
	UpdateTimeEmbed
	A struct{
		AA string `json:"aa" dbupdate:"aa"`
	}
	Test int
}

func TestStruct2Map(t *testing.T) {
	user := User{
		IdEmbed:         IdEmbed{ID: 7471066196823527424},
		Name:            "aavv",
		Nickname:		 "2233",
		TaskId:          8471066196823526789,
		UpdateTimeEmbed: UpdateTimeEmbed{UpdatedAt: time.Now()},
		Test: 19,
	}
	fmt.Println("Struct2MapByJson:", Struct2MapByJson(user))

	fmt.Println("Struct2Map dbupdate:", Struct2Map(user, "dbupdate"))
	fmt.Println("Struct2Map json:", Struct2Map(user))
}

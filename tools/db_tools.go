package tools

import (
	"github.com/bighuangbee/gokit/storage/kitGorm"
	"github.com/gogo/protobuf/proto"
	"time"
)


func PbToUpdateMap(message proto.Message, userId int64, userName string)(data map[string]interface{}) {
	data = PbToMapSkip(message)

	nowMyTime := kitGorm.MyTime(time.Now())
	data["updated_at"] = &nowMyTime
	data["updated_by"] = userId
	data["updated_by_name"] = userName

	return data
}

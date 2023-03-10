package tools

import (
	"bytes"
	"encoding/json"
	"github.com/bighuangbee/gokit/storage/kitGorm"
	"github.com/fatih/structs"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"time"
)


func PbToUpdateMap(message proto.Message, tableModel interface{}, userId int64, userName string)(data map[string]interface{}) {
	data = PbToModelMap(message, tableModel)

	nowMyTime := kitGorm.MyTime(time.Now())
	data["updated_at"] = &nowMyTime
	data["updated_by"] = userId
	data["updated_by_name"] = userName

	return data
}

//message: google.protobuf 自动填充为类型的零值，StringValue填充为空字符串、Int32Value填充为0，nil值会移除
//tableModel: db模型
func PbToModelMap(message proto.Message, tableModel interface{})(data map[string]interface{}){
	m := jsonpb.Marshaler{
		EnumsAsInts:  false,// 是否将枚举值设定为整数，而不是字符串类型
		EmitDefaults: true, // 是否将字段值为空的渲染到JSON结构中, nil被忽略，0或""保留
		OrigName:     true, // //是否使用原生的proto协议中的字段
	}
	var _buffer  bytes.Buffer
	err := m.Marshal(&_buffer, message)
	if err != nil{
		return
	}

	err = json.Unmarshal(_buffer.Bytes(), &data)
	if err != nil {
		return
	}

	//表字段
	tableFiles := make(map[string]interface{})
	s := structs.New(tableModel)
	s.TagName = "json"
	tableFiles = s.Map()

	dataFilter := make(map[string]interface{})
	for key, val := range data{
		if _, ok := tableFiles[key]; !ok{
			continue
		}
		m, ok := val.(map[string]interface{})
		if ok{
			if m != nil{
				if len(m) == 0{
					dataFilter[key] = ""
				}else{
					for _, v := range m {
						dataFilter[key] = v
					}
				}
			}
		}else{
			if val != nil{
				dataFilter[key] = val
			}
		}
	}
	return dataFilter
}

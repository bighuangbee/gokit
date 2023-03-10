package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"regexp"
)


func Struct2MapByJson(content interface{}) map[string]interface{} {
	var ret map[string]interface{}
	if marshalContent, err := json.Marshal(content); err == nil {
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		d.UseNumber() // 设置将float64转为number
		if err := d.Decode(&ret); err != nil {
			fmt.Println(err)
		} else {
			for k, v := range ret {
				ret[k] = v
			}
		}
	}else{
		fmt.Println(err)
	}
	return ret
}

//嵌套的字段被转换为map
func Struct2Map(data interface{}, tagName ...string)map[string]interface{}{
	ret := make(map[string]interface{})
	if data == nil{
		return ret
	}
	if len(tagName) == 0{
		tagName = append(tagName, "json")
	}
	s := structs.New(data)
	s.TagName = tagName[0]
	s.Name()
	ret = s.Map()

	for _, field := range s.Fields() {
		if field.IsEmbedded(){
			delete(ret, field.Name())
			for _, f := range s.Field(field.Name()).Fields() {
				if f.Tag(s.TagName) != ""{
					ret[f.Tag(s.TagName)] = f.Value()
				}
			}
		}

		//非json字段且tanName为空，移除
		if s.TagName != "json" && field.Tag(s.TagName) == ""{
			delete(ret, field.Name())
		}
	}
	return ret
}


func PbToMap(message proto.Message, filter ...string)(data map[string]interface{}){
	m := jsonpb.Marshaler{
		EnumsAsInts:  false,// 是否将枚举值设定为整数，而不是字符串类型
		EmitDefaults: false, // 是否将字段值为空的渲染到JSON结构中, nil被忽略，0或""保留
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

	dataFilter := make(map[string]interface{})
	if len(filter) > 0{
		for _, val := range filter{
			item, ok := data[val]
			if ok{
				dataFilter[val] = item
			}
		}
		return dataFilter
	}
	return
}


/**
 * @Description: 普通类型的字段设置了omitempty属性时，值为非空时保留key，值为空值(0/flase/"")忽略key
				 指针类型的字段设置了omitempty属性时，值为空值(0/flase/"")或非空值都保留key，值为nil忽略key
 * @param src 输入对象，struct类型或[]struct
 * @param target 转换目标 map/struct
 * @return error
*/
func StructTransform(src interface{}, target interface{}){
	tmp, err := json.Marshal(src)
	if err != nil{
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(tmp, &target)
	if err != nil{
		fmt.Println(err)
	}
}



//将字符串中的int64替换为string类型
func RegexpInt64(data string)(result string){
	reg := regexp.MustCompile(`id\":(\d{16,20}),"`)
	l := len(reg.FindAllString(data, -1)) //正则匹配16-20位的数字，如果找到就按正则替换并解析
	if l != 0 {
		result = reg.ReplaceAllString(data, `id": "${1}","`)
	}
	return result
}

//将struct中字段类型int64修改为string
func StructReviseInt64(obj interface{})(result interface{}){
	j, _ := json.Marshal(&obj)
	str := RegexpInt64(string(j))
	json.Unmarshal([]byte(str), &result)
	return
}


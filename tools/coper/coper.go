package coper

import (
"bytes"
"encoding/json"
	"github.com/petersunbag/coven"

	"github.com/golang/protobuf/jsonpb"
"github.com/golang/protobuf/proto"
)

// CopyToPBMessage Copy结构到proto.Message类型的结构
// 目前处理
// 	string => enum
//	time.Time => google.protobuf.Timestamp
func CopyToPBMessage(to proto.Message, fromList ...interface{}) error {
	unmarshaler := &jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
	for _, from := range fromList {
		b, _ := json.Marshal(from)
		if err := unmarshaler.Unmarshal(bytes.NewReader(b), to); err != nil {
			return err
		}
	}
	return nil
}

// Copy 普通的结构复制，必须想通类型
func Copy(to interface{}, fromList ...interface{}) error {
	for _, from := range fromList {
		b, _ := json.Marshal(from)
		if err := json.Unmarshal(b, to); err != nil {
			return err
		}
	}
	return nil
}

// CopyFromPBMessage Copy结构到proto.Message类型的结构
// 目前处理
//	Enum => string
func CopyFromPBMessage(to interface{}, from proto.Message) error {
	marshaler := jsonpb.Marshaler{EmitDefaults: true}
	s, _ := marshaler.MarshalToString(from)

	return json.Unmarshal([]byte(s), to)
}


func StructConvert(to, from interface{}) error {
	var c, err = coven.NewConverter(to, from)
	if err != nil {
		return err
	}
	err = c.Convert(to, from)
	if err != nil {
		return err
	}
	return nil
}

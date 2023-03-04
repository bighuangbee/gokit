package tools

import (
	"bytes"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// CopyToPBMessage Copy结构到proto.Message类型的结构
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

// Copy 普通的结构复制，必须相同类型
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


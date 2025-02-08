package function

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"time"
)

var TimeFormart = "2006-01-02 15:04:05"

func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func HmacSHA256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	result := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(result)
}

func InArray(str string, arr []string) bool {
	for _, val := range arr {
		if val == str {
			return true
		}
	}
	return false
}

func InArrayInt64(target int64, arr []int64) bool {
	for _, val := range arr {
		if val == target {
			return true
		}
	}
	return false
}

func SliceContains(src []int32, t int32) bool {
	for _, s := range src {
		if s == t {
			return true
		}
	}
	return false
}

func SliceContainsStr(src []string, t string) bool {
	for _, s := range src {
		if s == t {
			return true
		}
	}
	return false
}

func Date2timestamp(datetime string) (timestamp int64) {

	time, _ := time.ParseInLocation(TimeFormart, datetime, time.Local)
	timestamp = time.Unix()
	return timestamp
}

// 整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int32
	binary.Read(bytebuff, binary.LittleEndian, &data)
	return int(data)
}

package kitGorm

import (
	"regexp"
	"strings"
)

// IsUniqueErr 数据库唯一索引错误
func IsUniqueErr(err error) (is bool, key, value string) {
	msg := err.Error()
	reg, _ := regexp.Compile("Duplicate entry '.+' for key '.+'$")
	var strs []string
	str := reg.FindString(msg)
	if str != "" {
		reg, _ := regexp.Compile("'.+?'")
		strs = reg.FindAllString(msg, 2)
	}
	target := 2
	if len(strs) == target {
		for i, s := range strs {
			strs[i] = strings.Trim(s, "'")
		}

		return true, strs[1], strs[0]
	}
	return false, "", ""
}

// IsForeignKeyErr 数据库外键错误
func IsForeignKeyErr(err error) (is bool, foreignField, foreignTable, currentField string) {
	msg := err.Error()
	reg, _ := regexp.Compile("FOREIGN KEY .+ REFERENCES .+ .+$")
	var strList []string
	str := reg.FindString(msg)
	if str != "" {
		reg, _ = regexp.Compile("`.+?`")
		strList = reg.FindAllString(str, 3)
		target := 3
		if len(strList) == target {
			for i, s := range strList {
				strList[i] = strings.Trim(s, "`")
			}
			return true, strList[0], strList[1], strList[2]
		}
	}
	return false, "", "", ""
}

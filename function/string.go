package function

import (
	"math/rand"
	"strings"
)

func RandStr(n int) string {
	const letter = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+><?"
	letter_len := len(letter)

	b := make([]byte, n)
	for i := range b {
		b[i] = letter[rand.Intn(letter_len)]
	}
	return string(b)
}

/*
对切片的元素去重
*/
func SliceUnique(slice []string) []string {
	result := make([]string, 0, len(slice))
	temp := map[string]struct{}{}
	for _, item := range slice {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func SubStr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

// StrEndRepair
// @Description: 对字符串的小数位填充0，对齐长度
// @param str
// @param prec 对齐的小数位，小数点后不到prec位的填充0
// @return string
func StrEndRepair(str string, prec int) string {
	prec = prec + 1
	str = SubStr(str, 0, strings.Index(str, ".")+prec) //截取小数点后8位
	l := len(str[strings.Index(str, "."):])            //小数点后不够8位补0
	if l < prec {
		for i := 0; i < prec-l; i++ {
			str += "0"
		}
	}
	return str
}

/*
*
"23082383-62ac-4bde-8228-4ea734f74255, d1ca5687-d046-4bf9-b76e-672f2df8133b" =>
"'23082383-62ac-4bde-8228-4ea734f74255','d1ca5687-d046-4bf9-b76e-672f2df8133b'"
*/
func FormatWhereIn(str string) (whereIn string) {
	strSlice := strings.Split(str, ",")
	for key, val := range strSlice {
		whereIn += "'" + val + "'"
		if key < (len(strSlice) - 1) {
			whereIn += ","
		}
	}
	return
}

package strings2

import (
	"encoding/json"
	"strconv"
	"strings"
)

func IsNotBlank(str string) bool {
	return len(str) != 0
}

func Concat(items ...string) string {
	builder := strings.Builder{}

	for _, item := range items {
		builder.WriteString(item)
	}

	return builder.String()
}

func Join(items []string, sep string) string {
	return strings.Join(items, sep)
}

func Itoa(item int) string {
	return strconv.Itoa(item)
}

func ContainsIgnoreCase(strList []string, str string) bool {
	if len(strList) == 0 || len(str) == 0 {
		return false
	}

	for _, target := range strList {
		if EqualCaseIgnored(target, str) {
			return true
		}
	}

	return false
}

func ToJsonString(object interface{}) string {
	stringResult, err := json.Marshal(object)

	if err != nil {
		return `{"error": "failed to convert to json string"}`
	}

	return string(stringResult)
}

func EqualCaseIgnored(str1, str2 string) bool {
	if len(str1) == 0 && len(str2) == 0 {
		return true
	}

	if len(str1) == 0 || len(str2) == 0 {
		return false
	}

	return strings.ToLower(str1) == strings.ToLower(str2)
}

func Equals(str1, str2 string) bool {
	if len(str1) == 0 && len(str2) == 0 {
		return true
	}

	if len(str1) == 0 || len(str2) == 0 {
		return false
	}

	return str1 == str2
}

func NumericCompare(str1, str2 string) int {
	if len(str1) == 0 || len(str2) == 0 {
		return 0
	}

	if len(str1) == 0 {
		return -1
	}

	if len(str2) == 0 {
		return 1
	}

	float1, err1 := strconv.ParseFloat(str1, 64)
	float2, err2 := strconv.ParseFloat(str2, 64)

	if err1 != nil && err2 != nil {
		return 0
	}

	if err1 != nil {
		return -1
	}

	if err2 != nil {
		return 1
	}

	if float1 > float2 {
		return 1
	}

	if float1 == float2 {
		return 0
	}

	return -1
}

func Contains(src string, pattern string) bool {
	return strings.Contains(src, pattern)
}

func ToByte(src string) []byte {
	if IsNotBlank(src) {
		return []byte(src)
	}

	return nil
}

func Split(str string, separator string) []string {
	return strings.Split(str, separator)
}

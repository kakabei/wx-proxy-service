package common

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
	"wx-proxy-service/internal/types"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// LetterRunes 随机字符串字符池
var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func GenerateRandonString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = LetterRunes[rand.Intn(len(LetterRunes))]
	}
	return string(b)
}

func CreatRequestId() string {

	return fmt.Sprintf("wxmsg_%s", GenerateRandonString(10))
}

func GetRequstId(ctx context.Context) string {
	value := ctx.Value(types.CtxValRequestId{})
	if value != nil {
		if s, ok := value.(string); ok {
			return s
		}
	}

	return CreatRequestId()
}

// 支持两种格式  json格式：["xxx","yyy"], 分隔格式：xxx,yyy
func ParseStringToList(str string) ([]string, error) {
	var list []string

	if err := json.Unmarshal([]byte(str), &list); err == nil {
		return list, nil
	}

	list = strings.Split(str, ",")
	for i := range list {
		list[i] = strings.TrimSpace(list[i])
	}

	return list, nil
}

func ToJSON(object interface{}) string {
	bytes, _ := json.Marshal(object)
	return string(bytes)
}

// IsZero 判断val是否为Zero Value
func IsZero(val interface{}) bool {
	return reflect.Zero(reflect.TypeOf(val)).Interface() == val
}

// InArray 判断是否在数组中
func InArray(value int, arr []int) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}

// Uint64InArray 判断是否在数组中
func Uint64InArray(value uint64, arr []uint64) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}

// StringInArray 判断是否在数组中
func StringInArray(value string, arr []string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}

// RemoveFromStringArray 从数组中移除指定元素
func RemoveFromStringArray(arr []string, value ...string) []string {
	result := make([]string, 0)
	for _, v := range arr {
		if !StringInArray(v, value) {
			result = append(result, v)
		}
	}

	return result
}

// RemoveFromUint64Array 从数组中移除指定元素
func RemoveFromUint64Array(arr []uint64, value ...uint64) []uint64 {
	result := make([]uint64, 0)
	for _, v := range arr {
		if !Uint64InArray(v, value) {
			result = append(result, v)
		}
	}

	return result
}

func EqualIntArray(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

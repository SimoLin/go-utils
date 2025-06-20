package common

import (
	"cmp"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/elastic/elastic-agent-libs/mapstr"
)

// JSON字符串转Map[string]any
func StringToMap(s string) (M1 map[string]any, err error) {
	M1 = map[string]any{}
	err = json.Unmarshal([]byte(s), &M1)
	return
}

// JSON字符串转Map[string]any
func StringToSliceAny(s string) (S1 []any, err error) {
	S1 = []any{}
	err = json.Unmarshal([]byte(s), &S1)
	return
}

// JSON字符串转Map[string]any
func StringToSliceString(s string) (S1 []string, err error) {
	S1 = []string{}
	err = json.Unmarshal([]byte(s), &S1)
	return
}

// JSON字符串转Map[string]uint
func StringToSliceUint(s string) (S1 []uint, err error) {
	S1 = []uint{}
	err = json.Unmarshal([]byte(s), &S1)
	return
}

// 转JSON字符串
func MapToJsonString(M1 map[string]any) (r string) {
	result_byte, err := json.Marshal(M1)
	if err != nil {
		return
	}
	return string(result_byte)
}

// 获取 Map 的 key 值切片
func MapGetKeys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// 读取 map 中的值并返回指定类型的值
func MapGetValue[T any](M1 map[string]any, key string) (result T) {
	value, err := mapstr.M(M1).GetValue(key)
	if err != nil {
		return
	}
	if value, ok := value.(T); ok {
		return value
	}
	return result
}

// 读取 map 中的值并返回 string 类型的值，其他类型会强制使用 fmt.Sprintf("%v") 转换
func MapGetValueToString(M1 map[string]any, key string) (result string) {
	value, err := mapstr.M(M1).GetValue(key)
	if err != nil {
		return
	}
	if value, ok := value.(string); ok {
		return value
	}
	return fmt.Sprintf("%v", value)
}

// 转JSON字符串
func SliceToJsonString(S1 []any) (r string) {
	result_byte, err := json.Marshal(S1)
	if err != nil {
		return
	}
	return string(result_byte)
}

// 将任意类型的数组 S1 转为string字符串，间隔符 sep
func SliceToJoinString[T any](S1 []T, sep string) (result string) {
	switch len(S1) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%v", S1[0])
	}
	n := len(sep) * (len(S1) - 1)

	var b strings.Builder
	b.Grow(n)
	b.WriteString(fmt.Sprintf("%v", S1[0]))
	for _, s := range S1[1:] {
		b.WriteString(sep)
		b.WriteString(fmt.Sprintf("%v", s))
	}
	return b.String()
}

func SliceAnyToSliceString(S1 []any) (result []string) {
	result = []string{}
	switch len(S1) {
	case 0:
		return []string{}
	}
	for _, s := range S1 {
		result = append(result, fmt.Sprintf("%v", s))
	}
	return
}

// JSON反序列化获取的[]any，解析的数字格式是float64，转成[]uint
func SliceAnyToSliceUint(S1 []any) (result []uint) {
	result = []uint{}
	switch len(S1) {
	case 0:
		return []uint{}
	}
	for _, value := range S1 {
		value_float, ok := value.(float64)
		if !ok {
			continue
		}
		result = append(result, uint(value_float))
	}
	return
}

// 切片去重
func SliceUnique[T cmp.Ordered](S1 []T) []T {
	size := len(S1)
	if size == 0 {
		return []T{}
	}
	newSlices := make([]T, 0)
	m1 := make(map[T]byte)
	for _, v := range S1 {
		if _, ok := m1[v]; !ok {
			m1[v] = 1
			newSlices = append(newSlices, v)
		}
	}
	return newSlices
}

// 正则提取参数值(分组查询)
func RegexpSubmatch(str_input string, regex_expression string, index int) (r string) {
	re_match := regexp.MustCompile(regex_expression)
	result_match := re_match.FindStringSubmatch(str_input)
	if len(result_match) > index {
		r = result_match[index]
	}
	return
}

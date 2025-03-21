package common

import (
	"cmp"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/elastic/elastic-agent-libs/mapstr"
)

func MapToString(M1 map[string]any) (r string) {
	result_byte, err := json.Marshal(M1)
	if err != nil {
		return
	}
	return string(result_byte)
}

func StringToMap(s string) (M1 map[string]any, err error) {
	M1 = map[string]any{}
	err = json.Unmarshal([]byte(s), &M1)
	return
}

// 读取 map 中的值并返回指定类型的值
func MustGetValue[T any](M1 map[string]any, key string) (result T) {
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
func MustGetValueToString(M1 map[string]any, key string) (result string) {
	value, err := mapstr.M(M1).GetValue(key)
	if err != nil {
		return
	}
	if value, ok := value.(string); ok {
		return value
	}
	return fmt.Sprintf("%v", value)
}

// 将任意类型的数组 S1 转为string字符串，间隔符 sep
func MustGetStringJoin(S1 []any, sep string) (result string) {
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

// 正则提取参数值(分组查询)
func MustGetRegexpSubmatch(str_input string, regex_expression string, index int) (r string) {
	re_match := regexp.MustCompile(regex_expression)
	result_match := re_match.FindStringSubmatch(str_input)
	if len(result_match) > index {
		r = result_match[index]
	}
	return
}

// 获取 Map 的 key 值切片
func GetMapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// 切片去重
func GetUniqueSlice[T cmp.Ordered](S1 []T) []T {
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

package datetime

import (
	"fmt"
	"time"

	"github.com/jummyliu/pkg/datetime"
	"github.com/jummyliu/pkg/db/types"
)

// 获取 当前时间 标准格式 2006-01-02 15:04:05
func GetNowString() string {
	return datetime.FormatDate(time.Now())
}

// 获取 当前时间 简单格式 20060102-150405
func GetNowStringSimple() string {
	return datetime.FormatDateWithLayout(time.Now(), "20060102-150405")
}

// 基于 当前时间 给定 时间间隔 duration，返回时间字符串
func GetNowStringWithDuration(time_duration_string string) string {
	now_time := time.Now()
	time_duration, _ := time.ParseDuration(time_duration_string)
	result_time := now_time.Add(time_duration)
	return datetime.FormatDate(result_time)
}

func ParseDatetime(s string) types.Time {
	return types.Time(datetime.ParseDate(s))
}

// 基于 当前时间 给定 日期时间周期格式 返回日期时间范围 的字符串
//
//	datetime_type: 1周前, 2周前, 1个月前, 2个月前, 3个月前, 6个月前, 1年前
//	例子: GetDatetimeRange("1周前") => "2021-08-02", "2021-08-09", nil
func GetDatetimeRange(datetime_type string) (start_datetime string, end_datetime string, err error) {
	now_time := time.Now()
	end_datetime = datetime.FormatDate(now_time)
	result_time := now_time

	switch datetime_type {
	case "1周前":
		result_time = now_time.AddDate(0, 0, -7)
	case "2周前":
		result_time = now_time.AddDate(0, 0, -14)
	case "1个月前":
		result_time = now_time.AddDate(0, -1, 0)
	case "2个月前":
		result_time = now_time.AddDate(0, -2, 0)
	case "3个月前":
		result_time = now_time.AddDate(0, -3, 0)
	case "6个月前":
		result_time = now_time.AddDate(0, -6, 0)
	case "1年前":
		result_time = now_time.AddDate(-1, 0, 0)
	default:
		return "", "", fmt.Errorf("日期时间格式错误: %s", datetime_type)
	}

	start_datetime = datetime.FormatDate(result_time)

	return
}

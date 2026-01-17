package formatter

import (
	"github.com/spf13/cast"
)

// StrFormatter 转换str
type StrFormatter struct {
}

func (s StrFormatter) Name() string {
	return "str"
}

func (s StrFormatter) Format(value interface{}, param map[string]interface{}, args ...string) interface{} {
	return cast.ToString(value)
}

// StrToArrFormatter str换位数组
type StrToArrFormatter struct {
}

func (s StrToArrFormatter) Name() string {
	return "str_arr"
}

func (s StrToArrFormatter) Format(value interface{}, param map[string]interface{}, args ...string) interface{} {
	vv := cast.ToString(value)

	// 默认值
	split := ","
	t := "string"

	if len(args) > 0 && args[0] != "" {
		split = args[0]
	}
	if len(args) > 1 && args[1] != "" {
		t = args[1]
	}

	return SplitAndCast(vv, split, t)
}

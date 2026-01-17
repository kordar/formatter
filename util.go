package formatter

import (
	"strings"

	"github.com/spf13/cast"
)

func SplitAndCast(vv string, split string, t string) any {
	if vv == "" {
		// 根据类型返回空数组
		switch strings.ToLower(t) {
		case "int":
			return []int{}
		case "int64":
			return []int64{}
		case "float", "float64":
			return []float64{}
		case "bool":
			return []bool{}
		default:
			return []string{}
		}
	}

	items := strings.Split(vv, split)
	if len(items) == 0 {
		// 再次兜底空数组
		switch strings.ToLower(t) {
		case "int":
			return []int{}
		case "int64":
			return []int64{}
		case "float", "float64":
			return []float64{}
		case "bool":
			return []bool{}
		default:
			return []string{}
		}
	}

	switch strings.ToLower(t) {
	case "string":
		return items

	case "int":
		res := make([]int, 0, len(items))
		for _, v := range items {
			res = append(res, cast.ToInt(v))
		}
		return res

	case "int64":
		res := make([]int64, 0, len(items))
		for _, v := range items {
			res = append(res, cast.ToInt64(v))
		}
		return res

	case "float", "float64":
		res := make([]float64, 0, len(items))
		for _, v := range items {
			res = append(res, cast.ToFloat64(v))
		}
		return res

	case "bool":
		res := make([]bool, 0, len(items))
		for _, v := range items {
			res = append(res, cast.ToBool(v))
		}
		return res

	default:
		// 未知类型，返回 string
		return items
	}
}

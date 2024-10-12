package formatter

import (
	"fmt"
	"github.com/spf13/cast"
)

var container = map[string]Formatter{}

func init() {
	Register(StrFormatter{})
	Register(IntFormatter{})
	Register(LongFormatter{})
	Register(TimeFormatter{})
	Register(RoundStrFormatter{})
	Register(SerialFormatter{})
}

type Formatter interface {
	Name() string
	Format(value interface{}, param map[string]interface{}, args ...string) interface{}
}

func Register(formatters ...Formatter) {
	if formatters == nil {
		return
	}
	for _, formatter := range formatters {
		container[formatter.Name()] = formatter
	}
}

func Get(name string) Formatter {
	if formatter, ok := container[name]; ok {
		return formatter
	} else {
		return container["str"]
	}
}

func Format(name string, value interface{}, param map[string]interface{}, args ...string) interface{} {
	f := Get(name)
	if f == nil {
		return nil
	}
	return f.Format(value, param, args...)
}

// StrFormatter 转换str
type StrFormatter struct {
}

func (s StrFormatter) Name() string {
	return "str"
}

func (s StrFormatter) Format(value interface{}, param map[string]interface{}, args ...string) interface{} {
	return cast.ToString(value)
}

// IntFormatter 转int
type IntFormatter struct {
}

func (i IntFormatter) Name() string {
	return "int"
}

func (i IntFormatter) Format(value interface{}, param map[string]interface{}, args ...string) interface{} {
	return cast.ToInt(value)
}

// LongFormatter 转long
type LongFormatter struct {
}

func (i LongFormatter) Name() string {
	return "long"
}

func (i LongFormatter) Format(value interface{}, param map[string]interface{}, args ...string) interface{} {
	return cast.ToInt64(value)
}

// RoundStrFormatter 转long
type RoundStrFormatter struct {
}

func (i RoundStrFormatter) Name() string {
	return "roundstr"
}

func (i RoundStrFormatter) Format(value interface{}, param map[string]interface{}, args ...string) interface{} {
	format := "%.2f"
	if len(args) > 0 && args[0] != "" {
		d := cast.ToInt(args[0])
		format = fmt.Sprintf("%%.%df", d)
	}
	num := cast.ToFloat64(value)
	return fmt.Sprintf(format, num)
}

// TimeFormatter 转日期
type TimeFormatter struct {
}

func (i TimeFormatter) Name() string {
	return "time"
}

func (i TimeFormatter) Format(value interface{}, param map[string]interface{}, args ...string) interface{} {
	layout := "2006-01-02 15:04:05"
	if len(args) > 0 && args[0] != "" {
		layout = cast.ToString(args[0])
	}
	t := cast.ToTime(value)
	return t.Format(layout)
}

// SerialFormatter 转序号
type SerialFormatter struct {
}

func (s SerialFormatter) Name() string {
	return "serial"
}

func (s SerialFormatter) Format(value interface{}, param map[string]interface{}, args ...string) interface{} {
	if param == nil {
		return 0
	}
	page := cast.ToInt64(param["page"])
	pageSize := cast.ToInt64(param["pageSize"])
	index := cast.ToInt64(param["index"])
	return (page-1)*pageSize + index
}

package formatter

import (
	"fmt"
	"reflect"
	"strings"
)

var cached = map[string][]fieldCache{}

type fieldCache struct {
	Name      string
	FieldName string
	FieldType string
	Index     int
	Func      funcParam
}

type funcParam struct {
	Type string
	Args []string
}

func parseFormatterTag(value string, fieldType string) funcParam {
	if value == "" {
		return funcParam{Type: fieldType, Args: []string{}}
	}
	item := strings.Split(value, ":")
	param := funcParam{Type: item[0], Args: []string{}}
	if len(item) == 2 {
		args := strings.Split(item[1], ",")
		param.Args = append(param.Args, args...)
	}
	return param
}

func getCached(v interface{}) []fieldCache {
	t := reflect.TypeOf(v)
	p := t.PkgPath()
	key := fmt.Sprintf("%s.%s", p, t.Name())
	if cached[key] != nil {
		return cached[key]
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	m := make([]fieldCache, 0, t.NumField())
	// 遍历结构体字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)     // 获取第 i 个字段
		fieldName := field.Name // 字段名
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		name := fieldName
		if jsonTag != "" {
			name = jsonTag
		}
		// 获取字段标签
		tag := field.Tag.Get("formatter")
		f := parseFormatterTag(tag, field.Type.Name())
		mm := fieldCache{
			Name:      name,
			Func:      f,
			FieldName: fieldName,
			FieldType: field.Type.Name(),
			Index:     i,
		}
		m = append(m, mm)
	}
	return m
}

func ToDataList(v interface{}, params map[string]interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	if params == nil {
		params = map[string]interface{}{}
	}
	// 获取传入切片的值
	slice := reflect.ValueOf(v)
	if slice.Kind() == reflect.Ptr {
		slice = slice.Elem()
	}
	data := make([]map[string]interface{}, 0, slice.Len())
	// 将每个元素添加到 []interface{} 切片中
	for i := 0; i < slice.Len(); i++ {
		params["index"] = i
		vv := ToView(slice.Index(i).Interface(), params)
		data = append(data, vv)
	}
	return data
}

func ToView(v interface{}, params map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	if v == nil {
		return m
	}
	fields := getCached(v)
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	for _, field := range fields {
		f := Get(field.Func.Type)
		args := field.Func.Args
		val := value.Field(field.Index).Interface()
		key := field.Name
		m[key] = f.Format(val, params, args...)
	}
	return m
}

package formatter

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// 缓存结构体字段信息
var (
	cached   = make(map[string][]fieldCache)
	cacheMux sync.RWMutex
)

type fieldCache struct {
	Name      string // json 名称或字段名
	FieldName string // 结构体字段名
	FieldType string // 字段类型
	Index     int    // 字段索引
	Func      funcParam
}

type funcParam struct {
	Type string   // 格式化器类型
	Args []string // 格式化器参数
}

// 解析 formatter 标签，例如 "number:2,3"
func parseFormatterTag(value string, fieldType string) funcParam {
	if value == "" {
		return funcParam{Type: fieldType, Args: []string{}}
	}
	item := strings.SplitN(value, ":", 2)
	param := funcParam{Type: item[0], Args: []string{}}
	if len(item) == 2 {
		param.Args = strings.Split(item[1], ",")
	}
	return param
}

// 获取结构体字段缓存
func getCached(v interface{}) []fieldCache {
	t := reflect.TypeOf(v)
	if t == nil {
		return nil
	}

	// 获取实际类型，如果是指针则取 Elem
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	key := fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())

	// 使用读写锁保护缓存
	cacheMux.RLock()
	if fc, ok := cached[key]; ok {
		cacheMux.RUnlock()
		return fc
	}
	cacheMux.RUnlock()

	fields := make([]fieldCache, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 忽略私有字段和 json:"-"
		if field.PkgPath != "" {
			continue
		}
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		name := field.Name
		if jsonTag != "" {
			name = strings.Split(jsonTag, ",")[0]
		}

		tag := field.Tag.Get("formatter")
		f := parseFormatterTag(tag, field.Type.Name())

		fields = append(fields, fieldCache{
			Name:      name,
			FieldName: field.Name,
			FieldType: field.Type.Name(),
			Index:     i,
			Func:      f,
		})
	}

	cacheMux.Lock()
	cached[key] = fields
	cacheMux.Unlock()

	return fields
}

// 将切片转为 []map[string]interface{}
func ToDataList(v interface{}, params map[string]interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	if params == nil {
		params = make(map[string]interface{})
	}

	slice := reflect.ValueOf(v)
	if slice.Kind() == reflect.Ptr {
		slice = slice.Elem()
	}
	if slice.Kind() != reflect.Slice && slice.Kind() != reflect.Array {
		return nil
	}

	data := make([]map[string]interface{}, 0, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		params["index"] = i
		row := slice.Index(i).Interface()
		params["row"] = row // 当前整行传入 Format
		data = append(data, ToView(row, params))
	}
	return data
}

// 将单行结构体转为 map[string]interface{}
func ToView(v interface{}, params map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	if v == nil {
		return m
	}

	fields := getCached(v)
	if fields == nil {
		return m
	}

	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if params == nil {
		params = make(map[string]interface{})
	}
	params["row"] = v

	for _, field := range fields {
		f := Get(field.Func.Type) // 获取对应格式化器
		val := value.Field(field.Index).Interface()
		m[field.Name] = f.Format(val, params, field.Func.Args...)
	}
	return m
}

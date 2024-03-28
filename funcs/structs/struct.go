package structs

import (
	"reflect"
	"strings"
)

// GetTagJsonName 获取结构体json所有字段值，不包括`json:"-"`忽略的字段，英文逗号拼接
// withBackQuote 是否给字段值加反义字符(``)
func GetTagJsonName(s interface{}, withBackQuote bool) string {
	ptrType := reflect.TypeOf(s)
	if ptrType.Kind() != reflect.Ptr {
		return ""
	}
	strType := ptrType.Elem()
	fieldNum := strType.NumField()
	if fieldNum == 0 {
		return ""
	}
	jvs := make([]string, 0)
	for i := 0; i < fieldNum; i++ {
		jv := strType.Field(i).Tag.Get("json")
		if jv == "" || jv == "-" {
			continue
		}
		if i := strings.Index(jv, ","); i != -1 {
			jv = jv[0:i]
		}
		if withBackQuote {
			jvs = append(jvs, "`"+jv+"`")
		} else {
			jvs = append(jvs, jv)
		}
	}
	return strings.Join(jvs, ",")
}

package mysql

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func (m *Model) Session() *gorm.DB {
	return m.DB().Table(m.Table)
}

func (m *Model) GetPrimaryKeyFieldName() (string, error) {
	ptrType := reflect.TypeOf(m.Data)
	if ptrType.Kind() != reflect.Ptr {
		return "", ErrMustPtrData
	}
	strType := ptrType.Elem()
	fieldNum := strType.NumField()
	if fieldNum == 0 {
		return "", ErrMissPrimaryKey
	}
	for i := 0; i < fieldNum; i++ {
		if s := strType.Field(i).Tag.Get("gorm"); s != "" {
			if strings.Contains(strings.ToUpper(s), "PRIMARYKEY") {
				return strType.Field(i).Name, nil
			}
		}
	}
	return "", ErrMissPrimaryKey
}

func (m *Model) AllFields(withBackQuote bool) string {
	ptrType := reflect.TypeOf(m.Data)
	if ptrType.Kind() != reflect.Ptr {
		panic(ErrMustPtrData)
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

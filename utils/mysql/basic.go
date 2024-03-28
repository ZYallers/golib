package mysql

import (
	"gorm.io/gorm/schema"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func (m *Model) Session() *gorm.DB {
	return m.DB().Table(m.Table)
}

func (m *Model) GetPrimaryKeyFieldName() (string, error) {
	if m.Data == nil {
		return "", ErrNilPointer
	}

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
		if text := strType.Field(i).Tag.Get("gorm"); text != "" {
			if tagSetting := schema.ParseTagSetting(text, ";"); len(tagSetting) > 0 {
				if pk, ok := tagSetting["PRIMARYKEY"]; ok && pk == "PRIMARYKEY" {
					if col, ok := tagSetting["COLUMN"]; ok && col != "" {
						return col, nil
					}
				}
			}
		}
	}

	return "", ErrMissPrimaryKey
}

func (m *Model) AllFields(withBackQuote bool) string {
	if m.Data == nil {
		return ""
	}
	ptrType := reflect.TypeOf(m.Data)
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

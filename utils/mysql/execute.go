package mysql

import (
	"errors"
	"strings"
)

func (m *Model) Save(value interface{}, updates ...interface{}) (interface{}, error) {
	db := m.DB().Table(m.Table)
	if ul := len(updates); ul > 0 {
		if i, ok := updates[0].(int); ok && i > 0 {
			if ul > 1 {
				if s, ok := updates[1].(string); ok && s != "" {
					db = db.Select(strings.Split(s, ","))
				}
			}
			return value, db.Updates(value).Error
		}
	}
	return value, db.Create(value).Error
}

func (m *Model) SaveOrUpdate(value interface{}, primaryKey int, updateFields string) (interface{}, error) {
	db := m.DB().Table(m.Table)
	if primaryKey > 0 {
		if updateFields != "" {
			db = db.Select(strings.Split(updateFields, ","))
		}
		return value, db.Updates(value).Error
	}
	return value, db.Create(value).Error
}

func (m *Model) Update(where []interface{}, value interface{}) error {
	return m.DB().Table(m.Table).Where(where[0], where[1:]...).Updates(value).Error
}

func (m *Model) Delete(where []interface{}) error {
	if where == nil {
		return errors.New("query condition cannot be empty")
	}
	return m.DB().Table(m.Table).Where(where[0], where[1:]...).Delete(nil).Error
}

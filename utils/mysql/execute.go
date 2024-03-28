package mysql

import (
	"github.com/ZYallers/golib/funcs/conv"
	"gorm.io/gorm"
)

func (m *Model) Save(value interface{}, args ...interface{}) (interface{}, error) {
	db := m.Session()
	if argsLen := len(args); argsLen > 0 {
		if conv.ToInt(args[0]) > 0 {
			if argsLen > 1 && args[1] != nil && args[1] != "" {
				db = db.Select(args[1])
			}
			err := db.Updates(value).Error
			return value, err
		}
	}
	err := db.Create(value).Error
	return value, err
}

func (m *Model) SaveOrUpdate(value interface{}, primaryKey int, updateFields interface{}) (interface{}, error) {
	db := m.Session()
	if primaryKey > 0 {
		if updateFields != nil && updateFields != "" {
			db = db.Select(updateFields)
		}
		err := db.Updates(value).Error
		return value, err
	}
	err := db.Create(value).Error
	return value, err
}

func (m *Model) Create(value interface{}) (interface{}, error) {
	err := m.Session().Create(value).Error
	return value, err
}

func (m *Model) Update(where []interface{}, value interface{}) error {
	return m.Session().Where(where[0], where[1:]...).Updates(value).Error
}

func (m *Model) Delete(where []interface{}) error {
	if where == nil {
		return gorm.ErrMissingWhereClause
	}
	return m.Session().Where(where[0], where[1:]...).Delete(nil).Error
}

package mysql

import (
	"errors"
	"fmt"
	"github.com/ZYallers/golib/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"sync/atomic"
	"time"
)

type Model struct {
	Table string
	DB    func() *gorm.DB
}

const (
	defaultCharset          = "utf8mb4"
	defaultLoc              = "Local"
	defaultParseTime        = "true"
	defaultMaxAllowedPacket = "0"
	defaultTimeout          = "15s"
	retryMaxTimes           = 3
	defaultMaxIdleConns     = 100
	defaultMaxOpenConns     = 100
	defaultConnMaxLifetime  = 5 * time.Minute
)

func (m *Model) NewMysql(dbc *types.DBCollector, dialect *types.MysqlDialect) (*gorm.DB, error) {
	var err error
	for i := 1; i <= retryMaxTimes; i++ {
		if atomic.LoadUint32(&dbc.Done) == 0 {
			atomic.StoreUint32(&dbc.Done, 1)
			if dbc.Pointer, err = m.openMysql(dialect); err == nil && dbc.Pointer != nil {
				m.defaultConfig(dbc.Pointer)
			}
		}
		if err == nil {
			if dbc.Pointer == nil {
				err = fmt.Errorf("new mysql %s is nil", dialect.Db)
			} else {
				if sqlDB, err2 := dbc.Pointer.DB(); err2 != nil {
					err = err2
				} else {
					err = sqlDB.Ping()
				}
			}
		}
		if err != nil {
			atomic.StoreUint32(&dbc.Done, 0)
			if i < retryMaxTimes {
				time.Sleep(time.Duration(i*100) * time.Millisecond)
				continue
			} else {
				return nil, fmt.Errorf("new mysql %s error: %v", dialect.Db, err)
			}
		}
		break
	}
	return dbc.Pointer, nil
}

func (m *Model) openMysql(dialect *types.MysqlDialect) (*gorm.DB, error) {
	if dialect == nil {
		return nil, errors.New("mysql dialect is nil")
	}
	charset := defaultCharset
	if dialect.Charset != "" {
		charset = dialect.Charset
	}
	parseTime := defaultParseTime
	if dialect.ParseTime != "" {
		parseTime = dialect.ParseTime
	}
	loc := defaultLoc
	if dialect.Loc != "" {
		loc = dialect.Loc
	}
	maxAllowedPacket := defaultMaxAllowedPacket
	if dialect.MaxAllowedPacket != "" {
		maxAllowedPacket = dialect.MaxAllowedPacket
	}
	timeout := defaultTimeout
	if dialect.Timeout != "" {
		timeout = dialect.Timeout
	}
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s&maxAllowedPacket=%s&timeout=%s",
		dialect.User, dialect.Pwd, dialect.Host, dialect.Port, dialect.Db,
		charset, parseTime, loc, maxAllowedPacket, timeout)
	return gorm.Open(mysql.Open(dns), &gorm.Config{})
}

func (m *Model) defaultConfig(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		// 设置连接池中空闲连接的最大数量
		sqlDB.SetMaxIdleConns(defaultMaxIdleConns)
		// 设置打开数据库连接的最大数量
		sqlDB.SetMaxOpenConns(defaultMaxOpenConns)
		// 设置了连接可复用的最大时间
		sqlDB.SetConnMaxLifetime(defaultConnMaxLifetime)
	}
}

func (m *Model) Find(dest interface{}, where []interface{}, fields, order string, offset, limit int) {
	db := m.DB().Table(m.Table)
	if fields != "" {
		db = db.Select(fields)
	}
	if where != nil {
		db = db.Where(where[0], where[1:]...)
	}
	if order != "" {
		db = db.Order(order)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	db.Find(dest)
}

func (m *Model) FindOne(dest interface{}, where []interface{}, fields, order string) {
	m.Find(dest, where, fields, order, 0, 1)
}

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

func (m *Model) Count(where []interface{}) (count int64) {
	db := m.DB().Table(m.Table)
	if where != nil {
		db = db.Where(where[0], where[1:]...)
	}
	db.Count(&count)
	return
}

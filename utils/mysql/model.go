package mysql

import (
	"errors"
	"fmt"
	"github.com/ZYallers/golib/types"
	"gorm.io/gorm"
	"sync/atomic"
	"time"
)

const (
	retryMaxTimes          = 3
	retrySleepTime         = 100 * time.Millisecond
	defaultMaxIdleConns    = 100
	defaultMaxOpenConns    = 100
	defaultConnMaxLifetime = 5 * time.Minute
)

type Model struct {
	Table string
	DB    func() *gorm.DB
}

func (m *Model) NewMysql(dbc *types.DBCollector, dialect *types.MysqlDialect, config func() *gorm.Config) (*gorm.DB, error) {
	if dialect == nil {
		return nil, errors.New("mysql dialect is nil")
	}
	var err error
	for i := 1; i <= retryMaxTimes; i++ {
		if atomic.LoadUint32(&dbc.Done) == 0 {
			atomic.StoreUint32(&dbc.Done, 1)
			if dbc.Pointer, err = gorm.Open(m.Dialector(dialect), config()); err == nil && dbc.Pointer != nil {
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
				time.Sleep(retrySleepTime)
				continue
			} else {
				return nil, fmt.Errorf("new mysql %s error: %v", dialect.Db, err)
			}
		}
		break
	}
	return dbc.Pointer, nil
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

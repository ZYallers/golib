package mysql

import (
	"database/sql"
	"fmt"
	"github.com/ZYallers/golib/types"
	"gorm.io/gorm"
	"sync/atomic"
	"time"
)

const (
	retryMaxTimes          = 3
	retrySleepTime         = 100 * time.Millisecond
	defaultMaxIdleConns    = 25
	defaultMaxOpenConns    = 50
	defaultConnMaxLifetime = 5 * time.Minute
)

type Model struct {
	Table string
	DB    func() *gorm.DB
}

func (m *Model) NewMysql(dbc *types.DBCollector, mdt *types.MysqlDialect, cfg func() *gorm.Config) (*gorm.DB, error) {
	var err error
	for i := 1; i <= retryMaxTimes; i++ {
		if atomic.LoadUint32(&dbc.Done) == 0 {
			atomic.StoreUint32(&dbc.Done, 1)
			if dbc.Pointer, err = gorm.Open(m.Dialector(mdt), cfg()); err == nil {
				var db *sql.DB
				if db, err = dbc.Pointer.DB(); err == nil {
					db.SetMaxIdleConns(defaultMaxIdleConns)       // 设置连接池中空闲连接的最大数量
					db.SetMaxOpenConns(defaultMaxOpenConns)       // 设置打开数据库连接的最大数量
					db.SetConnMaxLifetime(defaultConnMaxLifetime) // 设置了连接可复用的最大时间
				}
			}
		} else {
			if dbc.Pointer == nil {
				err = fmt.Errorf("new mysql %s is nil", mdt.Db)
			} else {
				var db *sql.DB
				if db, err = dbc.Pointer.DB(); err == nil {
					err = db.Ping()
				}
			}
		}
		if err != nil {
			atomic.StoreUint32(&dbc.Done, 0)
			if i < retryMaxTimes {
				time.Sleep(retrySleepTime)
				continue
			} else {
				return nil, fmt.Errorf("new mysql %s error: %v", mdt.Db, err)
			}
		}
		break
	}
	return dbc.Pointer, nil
}

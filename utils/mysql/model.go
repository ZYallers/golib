package mysql

import (
	"database/sql"
	"fmt"
	"github.com/ZYallers/golib/types"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
)

const (
	retryMaxTimes          = 3
	retrySleepTime         = 100 * time.Millisecond
	defaultMaxIdleConns    = 5
	defaultMaxOpenConns    = 50
	defaultConnMaxIdleTime = 3 * time.Minute
	defaultConnMaxLifetime = 5 * time.Minute
)

type Model struct {
	Table string
	DB    func() *gorm.DB
}

func (m *Model) NewMysql(dbc *types.DBCollector, mdt *types.MysqlDialect, cfg func() *gorm.Config, opts ...interface{}) (*gorm.DB, error) {
	var err error
	for times := 1; times <= retryMaxTimes; times++ {
		dbc.Once(func() {
			if dbc.Pointer, err = gorm.Open(m.Dialector(mdt), cfg()); err == nil {
				var db *sql.DB
				if db, err = dbc.Pointer.DB(); err == nil {
					ol, maxIdle, maxOpen := len(opts), defaultMaxIdleConns, defaultMaxOpenConns
					maxIdleTime, maxLifeTime := defaultConnMaxIdleTime, defaultConnMaxLifetime
					if ol > 0 {
						maxIdle = cast.ToInt(opts[0])
					}
					if ol > 1 {
						maxOpen = cast.ToInt(opts[1])
					}
					if ol > 2 {
						maxIdleTime = cast.ToDuration(opts[2])
					}
					if ol > 3 {
						maxLifeTime = cast.ToDuration(opts[3])
					}
					db.SetMaxIdleConns(maxIdle)
					db.SetMaxOpenConns(maxOpen)
					db.SetConnMaxIdleTime(maxIdleTime)
					db.SetConnMaxLifetime(maxLifeTime)
				}
			}
		})

		if err == nil {
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
			if times < retryMaxTimes {
				dbc.Reset(func() { time.Sleep(retrySleepTime) })
				continue
			} else {
				return nil, fmt.Errorf("new mysql %s error: %v", mdt.Db, err)
			}
		}

		break
	}
	return dbc.Pointer, nil
}

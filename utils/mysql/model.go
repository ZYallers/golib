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
	retrySleepTime         = time.Second
	defaultMaxIdleConns    = 5
	defaultMaxOpenConns    = 50
	defaultConnMaxIdleTime = 5 * time.Minute
	defaultConnMaxLifetime = 10 * time.Minute
)

type Model struct {
	Table string
	DB    func() *gorm.DB
}

func (m *Model) NewMysql(dbc *types.DBCollector, mdt *types.MysqlDialect, f func() *gorm.Config, opts ...interface{}) (*gorm.DB, error) {
	var newErr error
	for i := 0; i < retryMaxTimes; i++ {
		dbc.Once(func() {
			cfg := f()
			cfg.DisableAutomaticPing = true
			if dbc.Pointer, newErr = gorm.Open(m.Dialector(mdt), cfg); newErr != nil {
				return
			}
			var db *sql.DB
			if db, newErr = dbc.Pointer.DB(); newErr != nil {
				return
			}
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
		})

		if newErr == nil {
			if dbc.Pointer == nil {
				newErr = fmt.Errorf("new mysql %s is nil", mdt.Db)
			} else {
				var db *sql.DB
				if db, newErr = dbc.Pointer.DB(); newErr == nil {
					newErr = db.Ping()
				}
			}
		}

		if newErr != nil {
			dbc.Reset(func() { time.Sleep(retrySleepTime) })
		} else {
			break
		}
	}
	return dbc.Pointer, newErr
}

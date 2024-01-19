package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
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

func (m *Model) NewMysql(collector *Collector, dialect *Dialect, f func() *gorm.Config, opts ...interface{}) (*gorm.DB, error) {
	var newErr error
	for i := 0; i < retryMaxTimes; i++ {
		collector.Once(func() {
			cfg := f()
			cfg.DisableAutomaticPing = true
			if collector.Pointer, newErr = gorm.Open(m.Dialector(dialect), cfg); newErr != nil {
				return
			}
			var db *sql.DB
			if db, newErr = collector.Pointer.DB(); newErr != nil {
				return
			}
			ol, maxIdle, maxOpen := len(opts), defaultMaxIdleConns, defaultMaxOpenConns
			maxIdleTime, maxLifeTime := defaultConnMaxIdleTime, defaultConnMaxLifetime
			if ol > 0 {
				maxIdle = opts[0].(int)
			}
			if ol > 1 {
				maxOpen = opts[1].(int)
			}
			if ol > 2 {
				maxIdleTime = opts[2].(time.Duration)
			}
			if ol > 3 {
				maxLifeTime = opts[3].(time.Duration)
			}
			db.SetMaxIdleConns(maxIdle)
			db.SetMaxOpenConns(maxOpen)
			db.SetConnMaxIdleTime(maxIdleTime)
			db.SetConnMaxLifetime(maxLifeTime)
		})

		if newErr == nil {
			if collector.Pointer == nil {
				newErr = fmt.Errorf("new mysql %s is nil", dialect.Db)
			} else {
				var db *sql.DB
				if db, newErr = collector.Pointer.DB(); newErr == nil {
					newErr = db.Ping()
				}
			}
		}

		if newErr != nil {
			collector.Reset(func() { time.Sleep(retrySleepTime) })
		} else {
			break
		}
	}
	return collector.Pointer, newErr
}

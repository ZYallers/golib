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
	DefaultMaxIdleConns    = 5
	DefaultMaxOpenConns    = 50
	DefaultConnMaxIdleTime = 5 * time.Minute
	DefaultConnMaxLifeTime = 10 * time.Minute
)

type Model struct {
	Table string
	Data  interface{}
	DB    func() *gorm.DB
}

func (m *Model) NewMysql(collector *Collector, dialect *Dialect, options ...Option) (*gorm.DB, error) {
	var newErr error
	for i := 0; i < retryMaxTimes; i++ {
		collector.Once(func() {
			config := &modelConfig{}
			for _, option := range options {
				option(config)
			}
			if config.logger == nil {
				config.logger = m.NewLogger(dialect.Db)
			}
			gormConfig := &gorm.Config{DisableAutomaticPing: true, Logger: config.logger}
			if collector.Pointer, newErr = gorm.Open(m.Dialector(dialect), gormConfig); newErr != nil {
				return
			}
			var sqlDB *sql.DB
			if sqlDB, newErr = collector.Pointer.DB(); newErr != nil {
				return
			}
			if config.maxIdleConns == 0 {
				config.maxIdleConns = DefaultMaxIdleConns
			}
			if config.maxOpenConns == 0 {
				config.maxOpenConns = DefaultMaxOpenConns
			}
			if config.connMaxIdleTime == 0 {
				config.connMaxIdleTime = DefaultConnMaxIdleTime
			}
			if config.connMaxLifeTime == 0 {
				config.connMaxLifeTime = DefaultConnMaxLifeTime
			}
			sqlDB.SetMaxIdleConns(config.maxIdleConns)
			sqlDB.SetMaxOpenConns(config.maxOpenConns)
			sqlDB.SetConnMaxIdleTime(config.connMaxIdleTime)
			sqlDB.SetConnMaxLifetime(config.connMaxLifeTime)
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

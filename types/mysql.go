package types

import (
	"sync"
	"sync/atomic"

	"gorm.io/gorm"
)

type MysqlDialect struct {
	User, Pwd, Host, Port, Db, Charset, Loc, ParseTime, MaxAllowedPacket, Timeout string
}

type DBCollector struct {
	done    uint32
	m       sync.Mutex
	Pointer *gorm.DB
}

func (d *DBCollector) Once(f func()) {
	if atomic.LoadUint32(&d.done) == 0 {
		d.doSlow(f)
	}
}

func (d *DBCollector) Reset(f func()) {
	d.m.Lock()
	defer d.m.Unlock()
	if d.done == 1 {
		defer atomic.StoreUint32(&d.done, 0)
		f()
	}
}

func (d *DBCollector) doSlow(f func()) {
	d.m.Lock()
	defer d.m.Unlock()
	if d.done == 0 {
		defer atomic.StoreUint32(&d.done, 1)
		f()
	}
}

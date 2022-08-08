package types

import (
	"gorm.io/gorm"
	"sync"
)

type MysqlDialect struct {
	User, Pwd, Host, Port, Db, Charset, Loc, ParseTime, MaxAllowedPacket, Timeout string
}

type DBCollector struct {
	Done    uint32
	M       sync.Mutex
	Pointer *gorm.DB
}

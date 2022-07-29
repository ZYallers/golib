package types

import "gorm.io/gorm"

type MysqlDialect struct {
	User, Pwd, Host, Port, Db, Charset, Loc, ParseTime, MaxAllowedPacket, Timeout string
}

type DBCollector struct {
	Done    uint32
	Pointer *gorm.DB
}

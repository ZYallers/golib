package mysql

import (
	"fmt"

	"github.com/ZYallers/golib/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DefaultCharset          = "utf8mb4"
	DefaultLoc              = "Local"
	DefaultParseTime        = "true"
	DefaultMaxAllowedPacket = "0"
	DefaultTimeout          = "15s"
)

func (m *Model) Dialector(dialect *types.MysqlDialect) gorm.Dialector {
	return mysql.Open(m.ParseDSN(dialect))
}

func (m *Model) ParseDSN(dialect *types.MysqlDialect) string {
	charset := DefaultCharset
	if dialect.Charset != "" {
		charset = dialect.Charset
	}
	parseTime := DefaultParseTime
	if dialect.ParseTime != "" {
		parseTime = dialect.ParseTime
	}
	loc := DefaultLoc
	if dialect.Loc != "" {
		loc = dialect.Loc
	}
	maxAllowedPacket := DefaultMaxAllowedPacket
	if dialect.MaxAllowedPacket != "" {
		maxAllowedPacket = dialect.MaxAllowedPacket
	}
	timeout := DefaultTimeout
	if dialect.Timeout != "" {
		timeout = dialect.Timeout
	}
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s&maxAllowedPacket=%s&timeout=%s",
		dialect.User, dialect.Pwd, dialect.Host, dialect.Port, dialect.Db,
		charset, parseTime, loc, maxAllowedPacket, timeout)
	return dns
}

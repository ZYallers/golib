package mysql

import (
	"fmt"
	"github.com/ZYallers/golib/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	defaultCharset          = "utf8mb4"
	defaultLoc              = "Local"
	defaultParseTime        = "true"
	defaultMaxAllowedPacket = "0"
	defaultTimeout          = "15s"
)

func (m *Model) Dialector(dialect *types.MysqlDialect) gorm.Dialector {
	return mysql.Open(m.ParseDSN(dialect))
}

func (m *Model) ParseDSN(dialect *types.MysqlDialect) string {
	charset := defaultCharset
	if dialect.Charset != "" {
		charset = dialect.Charset
	}
	parseTime := defaultParseTime
	if dialect.ParseTime != "" {
		parseTime = dialect.ParseTime
	}
	loc := defaultLoc
	if dialect.Loc != "" {
		loc = dialect.Loc
	}
	maxAllowedPacket := defaultMaxAllowedPacket
	if dialect.MaxAllowedPacket != "" {
		maxAllowedPacket = dialect.MaxAllowedPacket
	}
	timeout := defaultTimeout
	if dialect.Timeout != "" {
		timeout = dialect.Timeout
	}
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s&maxAllowedPacket=%s&timeout=%s",
		dialect.User, dialect.Pwd, dialect.Host, dialect.Port, dialect.Db,
		charset, parseTime, loc, maxAllowedPacket, timeout)
	return dns
}

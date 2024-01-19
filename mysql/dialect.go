package mysql

import (
	"fmt"

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

type Dialect struct {
	User             string
	Pwd              string
	Host             string
	Port             string
	Db               string
	Charset          string
	Loc              string
	ParseTime        string
	MaxAllowedPacket string
	Timeout          string
}

func (m *Model) Dialector(dialect *Dialect) gorm.Dialector {
	return mysql.Open(m.ParseDSN(dialect))
}

func (m *Model) ParseDSN(dialect *Dialect) string {
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

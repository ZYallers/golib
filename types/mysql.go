package define

type MysqlDialect struct {
	User, Pwd, Host, Port, Db, Charset, Loc, ParseTime, MaxAllowedPacket, Timeout string
}

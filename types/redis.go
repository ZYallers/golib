package define

type RedisClient struct {
	Host, Port, Pwd string
	Db              int
}

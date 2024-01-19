package trace

type Cache interface {
	Set(k, v []byte)
	Get([]byte) []byte
	Exist([]byte) bool
	Del([]byte)
	Clear()
}

package cache

type Cache interface {
	Get(key string) (data string, err error)
	Set(key, value string) error
}

package tinycache

type Cache interface {
	Set(key string, value string) (err error)
	Get(key string) (value string, err error)
	Del(key string) (err error)
	Exists(key string) (ret bool, err error)
	Total() (total int)
}

type Limit struct {
	TotalElements int
	SizeInBytes   int
}

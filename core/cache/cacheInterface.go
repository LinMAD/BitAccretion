package cache

// ICache represents basic method for cache handler
type ICache interface {
	Add(key string, data interface{})
	Get(key string) interface{}
	DeleteBy(key string) bool
}

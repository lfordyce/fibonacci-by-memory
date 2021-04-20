package store

type KV interface {
	Set(k int64, v interface{}) error
	Get(k int64, v interface{}) (found bool, err error)
	Delete(k int64) error
	Close() error
}




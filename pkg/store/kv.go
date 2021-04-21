package store

import (
	"math/big"
)

type KV interface {
	Fibonacci(int) (*big.Int, error)
	Records(int) (int64, error)
	Truncate() (int64, error)
	Close() error
}
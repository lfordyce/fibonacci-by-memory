package store

import (
	"math/big"
)

// KV is an abstraction for a fibonacci key-value store.
type KV interface {
	// Fibonacci calculate the nth value in the fibonacci sequence.
	Fibonacci(int) (*big.Int, error)
	// Records finds the number of memoized results less than a given value.
	Records(int) (int64, error)
	// Truncate clears all the memorized results.
	Truncate() (int64, error)
	// Close will close the connection to the DB server.
	Close() error
}

package postgres

import (
	gosql "database/sql"
	"github.com/lfordyce/fibonacci-by-memory/pkg/store/sql"
	"github.com/lib/pq"
)

// Client is a KV implementation for PostgreSQL.
type Client struct {
	*sql.Client
}

// Option is an abstraction for providing addition configuration
// to the database connection.
type Option interface {
	Apply(*gosql.DB)
}

// OptionFunc is the concrete type to implement the Option interface
type OptionFunc func(*gosql.DB)

// Apply satisfies the Option interface
func (fn OptionFunc) Apply(db *gosql.DB) {
	fn(db)
}

// WithMaxOpenConns limit number of concurrent connections. Typical max connections on a PostgreSQL server is 100.
// This prevents "Error 1040: Too many connections", which otherwise occurs for example with 500 concurrent goroutines.
func WithMaxOpenConns(n int) Option {
	return OptionFunc(func(db *gosql.DB) {
		db.SetMaxOpenConns(n)
	})
}

// NewClient creates a new PostgreSQL client.
//
// Close() method must be called on the client when you're done working with it.
func NewClient(addr string, opts ...Option) (Client, error) {
	pqURL, err := pq.ParseURL(addr)
	if err != nil {
		return Client{}, err
	}
	conn, err := gosql.Open("postgres", pqURL)
	if err != nil {
		return Client{}, err
	}
	if err := conn.Ping(); err != nil {
		return Client{}, err
	}

	for _, opt := range opts {
		opt.Apply(conn)
	}

	fibonacciStmt, err := conn.Prepare("SELECT fibonacci_cached($1)")
	if err != nil {
		return Client{}, nil
	}

	truncateStmt, err := conn.Prepare("TRUNCATE fib_store")
	if err != nil {
		return Client{}, err
	}

	resultCountStmt, err := conn.Prepare("SELECT count(fibonacci) FROM fib_store WHERE fibonacci < $1")
	if err != nil {
		return Client{}, err
	}

	c := sql.Client{
		C:             conn,
		FibonacciStmt: fibonacciStmt,
		TruncateStmt:  truncateStmt,
		CountStmt:     resultCountStmt,
	}
	return Client{&c}, nil
}

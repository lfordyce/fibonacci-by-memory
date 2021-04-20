package postgres

import (
	gosql "database/sql"
	"github.com/lfordyce/fibonacci-by-memory/pkg/store"
	"github.com/lfordyce/fibonacci-by-memory/pkg/store/sql"
	"github.com/lib/pq"
)

// Client is a KV implementation for PostgreSQL.
type pqClient struct {
	*sql.Client
}

func NewClient(addr string) (store.KV, error) {
	pqURL, err := pq.ParseURL(addr)
	if err != nil {
		return pqClient{}, err
	}
	conn, err := gosql.Open("postgres", pqURL)
	if err != nil {
		return pqClient{}, err
	}
	if err := conn.Ping(); err != nil {
		return pqClient{}, err
	}

	// Limit number of concurrent connections. Typical max connections on a PostgreSQL server is 100.
	// This prevents "Error 1040: Too many connections", which otherwise occurs for example with 500 concurrent goroutines.
	conn.SetMaxOpenConns(10)

	upsertStmt, err := conn.Prepare("INSERT INTO fib_cache (ordinal, fibonacci) VALUES ($1, $2) ON CONFLICT (ordinal) DO UPDATE SET fibonacci = $2")
	if err != nil {
		return pqClient{}, err
	}

	getStmt, err := conn.Prepare("SELECT v FROM fib_cache  WHERE ordinal = $1")
	if err != nil {
		return pqClient{}, err
	}

	deleteStmt, err := conn.Prepare("DELETE FROM fib_cache WHERE ordinal = $1")
	if err != nil {
		return pqClient{}, err
	}

	c := sql.Client{
		C:          conn,
		UpsertStmt: upsertStmt,
		GetStmt:    getStmt,
		DeleteStmt: deleteStmt,
	}
	return pqClient{&c}, nil
}

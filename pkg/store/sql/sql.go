package sql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"math/big"
)

type Client struct {
	C             *sql.DB
	FibonacciStmt *sql.Stmt
	TruncateStmt  *sql.Stmt
	CountStmt     *sql.Stmt
}

type BigInt struct {
	big.Int
}

func (b *BigInt) Value() (driver.Value, error) {
	if b != nil {
		return b.String(), nil
	}
	return nil, nil
}

func (b *BigInt) Scan(value interface{}) error {
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	if _, ok := b.SetString(i.String, 10); ok {
		return nil
	}
	return fmt.Errorf("could not scan type %T into BigInt", value)
}

func (c Client) Records(n int) (int64, error) {
	var resultSet int64
	if err := c.CountStmt.QueryRow(n).Scan(&resultSet); err != nil {
		return 0, err
	}
	return resultSet, nil
}

func (c Client) Truncate() (int64, error) {
	result, err := c.TruncateStmt.Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (c Client) Close() error {
	return c.C.Close()
}

func (c Client) Fibonacci(n int) (*big.Int, error) {
	var bigint BigInt
	if err := c.FibonacciStmt.QueryRow(n).Scan(&bigint); err != nil {
		return nil, err
	}
	return &bigint.Int, nil
}

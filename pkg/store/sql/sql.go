package sql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"math/big"
)

type Client struct {
	C          *sql.DB
	UpsertStmt *sql.Stmt
	GetStmt    *sql.Stmt
	DeleteStmt *sql.Stmt
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

func (c Client) Set(k int64, v interface{}) error {
	_, err := c.UpsertStmt.Exec(k, v)
	if err != nil {
		return err
	}
	return err
}

func (c Client) Get(k int64, v interface{}) (found bool, err error) {
	// TODO: Consider using RawBytes.
	dataPtr := new([]byte)
	if err := c.GetStmt.QueryRow(k).Scan(dataPtr); err != nil {
		return false, err
	}
	// TODO
	return true, nil
}

func (c Client) Delete(k int64) error {
	if _, err := c.DeleteStmt.Exec(k); err != nil {
		return err
	}
	return nil
}

func (c Client) Close() error {
	return c.C.Close()
}

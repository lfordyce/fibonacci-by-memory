package store

import (
	"database/sql"
	"fmt"
	"github.com/lfordyce/fibonacci-by-memory/pkg/store/postgres"
	"github.com/lib/pq"
	"testing"
)

func TestFibonacciIntegration(t *testing.T) {
	connectionURL := "postgresql://postgres:changeme@localhost:5432/postgres?sslmode=disable"
	if !checkConnection(connectionURL, t) {
		t.Skip("No connection to PostgreSQL could be established. Probably not running in a proper test environment.")
	}

	client, err := postgres.NewClient(connectionURL)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	fibonacci, err := client.Fibonacci(10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fibonacci)
}

func checkConnection(connectionURL string, t *testing.T) bool {
	url, err := pq.ParseURL(connectionURL)
	if err != nil {
		t.Logf("An error occurred parsing the connection URL: %v\n", err)
		return false
	}
	db, err := sql.Open("postgres", url)
	if err != nil {
		t.Logf("An error occurred during testing the connection to the server: %v\n", err)
		return false
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Logf("An error occurred during testing the connection to the server: %v\n", err)
		return false
	}
	return true
}

package store

import (
	"fmt"
	"github.com/lfordyce/fibonacci-by-memory/pkg/store/postgres"
	"testing"
)

func TestFibonacciIntegration(t *testing.T) {

	// postgres://postgres:changeme@localhost:5432/postgres?sslmode=disable
	client, err := postgres.NewClient("postgresql://postgres:changeme@localhost:5432/postgres?sslmode=disable")
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
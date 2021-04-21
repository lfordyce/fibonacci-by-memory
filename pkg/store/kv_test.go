package store

import (
	"fmt"
	"github.com/lfordyce/fibonacci-by-memory/pkg/store/postgres"
	"github.com/ory/dockertest/v3"
	"testing"
)

func TestFibonacciIntegration(t *testing.T) {
	var client postgres.Client
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatal(err)
	}
	resource, err := pool.BuildAndRun("fibonacci-image", "../../psql/Dockerfile", []string{
		"POSTGRES_PASSWORD=changeme", "POSTGRES_DB=postgres",
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = pool.Retry(func() error {
		var err error
		addr := fmt.Sprintf("postgres://postgres:changeme@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), "postgres")
		client, err = postgres.NewClient(addr)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	result, err := client.Fibonacci(11)
	if err != nil {
		t.Fatal(err)
	}
	if actual := result.Int64(); actual != 89 {
		t.Errorf("expected the (%d)th value in the fibonacci sequence to be: %d, actual %d", 11, 89, actual)
	}

	if err := client.Close(); err != nil {
		t.Errorf("failed to close connection to DB: %v", err)
	}
	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		t.Fatal(err)
	}
}

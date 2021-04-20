package cmd

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lfordyce/fibonacci-by-memory/log"
	"github.com/lfordyce/fibonacci-by-memory/pkg/store"
	"github.com/lfordyce/fibonacci-by-memory/pkg/store/postgres"
	"github.com/urfave/cli/v2"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

type config struct {
	Addr string
	// PostgreSQL Connection URL.
	ConnectionURL string
	// Limits the number of open connections to the PostgreSQL server.
	// -1 for no limit. 0 will lead to the default value (100) being set.
	// Optional (100 by default).
	MaxOpenConnections int
}

func (cfg *config) cliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "bind", Value: "0.0.0.0:8000", Destination: &cfg.Addr},
		&cli.StringFlag{Name: "postgres", Value: "postgres://postgres:changeme@localhost:5432/postgres?sslmode=disable", Destination: &cfg.ConnectionURL},
		&cli.IntFlag{Name: "max_connections", Value: 100, Destination: &cfg.MaxOpenConnections},
	}
}

type Server struct {
	*http.Server
	log.Logger
}

func newServer(addr string, h http.Handler, l log.Logger) Server {
	s := Server{
		Server: &http.Server{
			Addr:         addr,
			Handler:      h,
			WriteTimeout: 60 * time.Second,
			ReadTimeout:  60 * time.Second,
			ErrorLog:     l.StandardLog("fibonacci.http.server"),
		}, Logger: l}
	return s
}

type FibServer interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func Command(ac chan FibServer) *cli.Command {
	cfg := new(config)
	return &cli.Command{
		Name:  "server",
		Flags: append([]cli.Flag{}, cfg.cliFlags()...),
		Action: func(c *cli.Context) error {

			dbClient, err := postgres.NewClient(cfg.ConnectionURL)
			if err != nil {
				return err
			}

			l := log.NewLogger(log.WithApp())

			mux := chi.NewRouter()
			mux.Route("/v1", func(r chi.Router) {
				r.Mount("/api", MakeFibRoute(dbClient))
			})

			ac <- newServer(cfg.Addr, mux, l)
			return nil
		},
	}
}

func MakeFibRoute(store store.KV) http.Handler {
	router := chi.NewRouter()
	router.Route("/fib/", func(r chi.Router) {
		r.Get("/{ordinal}", func(w http.ResponseWriter, r *http.Request) {
			ordinal := chi.URLParam(r, "ordinal")

			cache := make(map[int64]int64)
			cache[1] = 1
			cache[2] = 1
			newFib := memoize(fib, cache)

			parseOrdinal, err := strconv.ParseInt(ordinal, 10, 64)
			if err != nil {
				http.Error(w, fmt.Errorf("ordinal value must be a number: %w", err).Error(), http.StatusBadRequest)
				return
			}
			i := newFib(parseOrdinal)
			if _, err := w.Write([]byte(strconv.Itoa(int(i)))); err != nil {
				http.Error(w, fmt.Errorf("failed to convert fibonicci result to string: %w", err).Error(), http.StatusInternalServerError)
				return
			}
			return
		})
	})
	return router
}

func fib(n int64) int64 {
	if n == 1 || n == 2 {
		return 1
	}
	return fib(n-2) + fib(n-1)
}

func memoize(targetFunc func(int64) int64, cache map[int64]int64) func(int64) int64 {
	middleLayer := func(n int64) int64 {
		if cache[n] != 0 {
			return cache[n]
		}
		return targetFunc(n)
	}
	return func(n int64) int64 {
		for cache[n] == 0 {
			cache[n] = middleLayer(n-1) + middleLayer(n-2)
		}
		return cache[n]
	}
}

func bigFib(n int64) *big.Int {
	switch n {
	case 0:
		return big.NewInt(0)
	case 1:
		return big.NewInt(1)
	case 2:
		return big.NewInt(1)
	default:
		a := bigFib(n - 2)
		b := bigFib(n - 1)
		value := big.NewInt(0)
		value.Add(value, a)
		value.Add(value, b)

		return value
	}
}

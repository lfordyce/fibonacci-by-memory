package cmd

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lfordyce/fibonacci-by-memory/log"
	"github.com/lfordyce/fibonacci-by-memory/pkg/store"
	"github.com/lfordyce/fibonacci-by-memory/pkg/store/postgres"
	"github.com/urfave/cli/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type config struct {
	// Server connection address.
	// default value: 0.0.0.0:8000
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

			dbClient, err := postgres.NewClient(cfg.ConnectionURL, postgres.WithMaxOpenConns(cfg.MaxOpenConnections))
			if err != nil {
				return err
			}

			l := log.NewLogger(log.WithApp())

			mux := chi.NewRouter()
			mux.Use(l.LoggerMiddleware())
			mux.Route("/v1", func(r chi.Router) {
				r.Mount("/api", MakeFibonacciRoute(dbClient))
			})

			logRoutes(mux, l)

			ac <- newServer(cfg.Addr, mux, l)
			return nil
		},
	}
}

func MakeFibonacciRoute(store store.KV) http.Handler {
	router := chi.NewRouter()
	router.Route("/fib/", func(r chi.Router) {
		// fetch the Fibonacci number given an ordinal
		r.Get("/{ordinal}", func(w http.ResponseWriter, r *http.Request) {
			ordinal := chi.URLParam(r, "ordinal")

			parseOrdinal, err := strconv.ParseInt(ordinal, 10, 64)
			if err != nil {
				http.Error(w, fmt.Errorf("ordinal value must be a number: %w", err).Error(), http.StatusBadRequest)
				return
			}

			fibonacci, err := store.Fibonacci(int(parseOrdinal))
			if err != nil {
				errResp := fmt.Errorf("failed to calculate the fibonacci value at ordinal %d: %w", parseOrdinal, err)
				http.Error(w, errResp.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := w.Write([]byte(fibonacci.String())); err != nil {
				http.Error(w, fmt.Errorf("failed to write fibonicci result response: %w", err).Error(), http.StatusInternalServerError)
				return
			}
			return
		})
		// fetch the number of memoized results less than a given value (e.g. there are 12 intermediate results less than 120)
		r.Get("/records/{count}", func(w http.ResponseWriter, r *http.Request) {
			count := chi.URLParam(r, "count")

			parseCount, err := strconv.ParseInt(count, 10, 64)
			if err != nil {
				http.Error(w, fmt.Errorf("count value must be a number: %w", err).Error(), http.StatusBadRequest)
				return
			}
			records, err := store.Records(int(parseCount))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := w.Write([]byte(strconv.Itoa(int(records)))); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		})
		// clear the data store
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			rowsAffected, err := store.Truncate()
			if err != nil {
				http.Error(w, fmt.Errorf("failed to clear fib_store records %w", err).Error(), http.StatusInternalServerError)
				return
			}
			if _, err := w.Write([]byte(fmt.Sprintf("total rows deleted from fib_store: %d", rowsAffected))); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		})
	})
	return router
}

func logRoutes(mux *chi.Mux, l log.Logger) {
	if err := chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)

		l.Log("route.walk").Str("method", method).
			Str("route", route).Msg("")

		return nil
	}); err != nil {
		l.Err(err).Msg("Failed to walk routes")
	}
}

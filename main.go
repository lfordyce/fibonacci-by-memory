package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/lfordyce/fibonacci-by-memory/cmd"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
)

var (
	// Version contains the current version.
	Version = "dev"
	// BuildDate contains a string with the build date.
	//BuildDate = "unknown"
)

func main () {
	ac := make(chan cmd.FibServer, 1)

	app := cli.NewApp()
	app.Name = "fibmemo"
	app.Version = Version
	app.Commands = cli.Commands{
		cmd.Command(ac),
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}

	ps := <-ac

	closeReq := cmd.RegisterInterrupt()
	go func() {
		<-closeReq
		if err := ps.Shutdown(context.Background()); err != nil {
			panic(fmt.Errorf("could not shutdown server: %w", err))
		}
	}()

	if err := ps.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

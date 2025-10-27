package main

import (
	"context"
	"os"

	"github.com/anttiharju/relcheck/internal/buildinfo"
	"github.com/anttiharju/relcheck/internal/cli"
	"github.com/anttiharju/relcheck/internal/exitcode"
	"github.com/anttiharju/relcheck/internal/interrupt"
)

var (
	revision string
	version  string
	time     string
)

func main() {
	go interrupt.Listen(exitcode.Interrupt, os.Interrupt)

	ctx := context.Background()
	exitCode := cli.Start(ctx, buildinfo.New(revision, version, time), os.Args[1:])
	os.Exit(int(exitCode))
}

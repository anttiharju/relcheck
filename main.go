package main

import (
	"context"
	"os"

	"github.com/anttiharju/relcheck/internal/cli"
	"github.com/anttiharju/relcheck/internal/exitcode"
	"github.com/anttiharju/relcheck/pkg/interrupt"
)

func main() {
	go interrupt.Listen("relcheck", exitcode.Interrupt, os.Interrupt)

	ctx := context.Background()
	exitCode := cli.Start(ctx, os.Args[1:])
	os.Exit(exitCode)
}

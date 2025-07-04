package main

import (
	"context"
	"os"

	"github.com/anttiharju/relcheck/internal/cli"
	"github.com/anttiharju/relcheck/internal/exitcode"
	"github.com/anttiharju/relcheck/internal/interrupt"
)

func main() {
	go interrupt.Listen(exitcode.Interrupt, os.Interrupt)

	ctx := context.Background()
	exitCode := cli.Run(ctx, os.Args[1:])
	os.Exit(int(exitCode))
}

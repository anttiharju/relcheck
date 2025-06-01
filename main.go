package main

import (
	"context"
	"os"

	"github.com/anttiharju/relcheck/internal/exitcode"
	"github.com/anttiharju/relcheck/internal/program"
	"github.com/anttiharju/relcheck/pkg/interrupt"
)

func main() {
	go interrupt.Listen(exitcode.Interrupt, os.Interrupt)

	ctx := context.Background()
	exitCode := program.Start(ctx, os.Args[1:])
	os.Exit(exitCode)
}

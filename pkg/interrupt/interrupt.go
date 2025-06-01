package interrupt

import (
	"fmt"
	"os"
	"os/signal"
)

func Listen(name string, exitcode int, signals ...os.Signal) {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, signals...)
	<-interruptCh
	fmt.Printf("\n%s: interrupted\n", name) // leading \n to have ^C appear on its own line
	os.Exit(exitcode)
}

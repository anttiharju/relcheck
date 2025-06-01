package usage

import (
	"fmt"

	"github.com/anttiharju/relcheck/internal/exitcode"
)

func Print() exitcode.Exitcode {
	fmt.Println("Usage: relcheck [--verbose] [--color=always] <file1.md> [file2.md] ...")
	fmt.Println("   or: relcheck [--verbose] [--color=always] run  (to check all *.md files in Git)")
	fmt.Println("   or: relcheck version  (to show version information)")

	return exitcode.UsageError
}

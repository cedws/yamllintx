package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cedws/yamllintx/lint"
)

func main() {
	linters := lint.Chain{
		lint.Comments{
			RequireStartingSpace: true,
			IgnoreShebangs:       true,
		},
		lint.Octal{},
		lint.Hyphens{
			MaxSpacesAfter: 1,
		},
		lint.TrailingSpaces{},
		lint.Anchors(lint.AnchorOpts{
			ForbidUndeclaredAliases: true,
			ForbidDuplicatedAnchors: true,
			ForbidUnusedAnchors:     true,
		}),
		lint.Braces{
			Forbid:          lint.ForbidBracesNone,
			MinSpacesInside: 1,
			MaxSpacesInside: 3,
		},
	}

	file := flag.String("src", "", "source file")
	flag.Parse()

	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}

	bytes, err := os.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(os.Stderr, filepath.Clean(*file))

	for err := range lint.LintAll(bytes, linters...) {
		fmt.Fprintf(os.Stderr, "  %d:%d\t%s\n", err.Line, err.Column, err.Error)
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/itsliamegan/watch"
)

func main() {
	flag.Usage = func() {
		fmt.Println("watch - recursively poll a directory for file changes")
		fmt.Println()
		fmt.Println("usage: watch [<directory>]")
		flag.PrintDefaults()
	}
	flag.Parse()

	workingDir, err := os.Getwd()
	exitIf(err)

	var rootDir string
	if len(flag.Args()) > 0 {
		rootDir = flag.Arg(0)

		if !filepath.IsAbs(rootDir) {
			rootDir = filepath.Join(workingDir, rootDir)
		}
	} else {
		rootDir = workingDir
	}

	changes, errs := watch.Start(rootDir)
	for {
		select {
		case change := <-changes:
			fmt.Println(change)
		case err := <-errs:
			exitWith(err)
		}
	}
}

func exitIf(err error) {
	if err != nil {
		exitWith(err)
	}
}

func exitWith(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

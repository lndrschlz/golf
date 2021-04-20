package main

import (
	"os"
)

const (
	version = "0.8 [wip]"
)

// golf is a line facilitator which acts like preprocessor based on the go language syntax.
func main() {
	args := os.Args
	if len(args) == 1 {
		printTitle()
		printHelp()
		os.Exit(1)
	}
	source, dest, initialize, verbose := processArguments(args[1:])
	if verbose {
		printTitle()
	}
	checkFile(source)
	content := processFile(source, initialize, verbose)
	saveFile(dest, content)
}

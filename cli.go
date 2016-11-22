package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

func (cli *CLI) Fatalf(format string, a ...interface{}) int {
	fmt.Fprintf(cli.errStream, "%s: %s\n", Name, fmt.Sprintf(format, a...))
	return ExitCodeError
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		prefix string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&prefix, "prefix", "release-branch.", "Release branch prefix")
	flags.StringVar(&prefix, "p", "release-branch.", "Release branch prefix(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	err := flags.Parse(args[1:])
	if err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	var pkgroot string
	if len(os.Args) > 1 {
		pkgroot = os.Args[1]
	} else {
		pkgroot, err = os.Getwd()
		if err != nil {
			return ExitCodeError
		}
	}

	v, ret := findversion(cli, pkgroot, prefix)
	fmt.Fprintf(cli.outStream, "%s\n", v)

	return ret
}

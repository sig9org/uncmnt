package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/sig9org/uncmnt/internal/uncomment"
	"github.com/sig9org/uncmnt/internal/update"
	"github.com/sig9org/uncmnt/internal/version"
)

const toolName = "uncmnt"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printHelp(stdout)
		return 0
	}

	var showHelp, showVersion, doUpdate bool
	fs := flag.NewFlagSet(toolName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&showHelp, "h", false, "")
	fs.BoolVar(&showHelp, "help", false, "")
	fs.BoolVar(&showVersion, "v", false, "")
	fs.BoolVar(&showVersion, "version", false, "")
	fs.BoolVar(&doUpdate, "u", false, "")
	fs.BoolVar(&doUpdate, "update", false, "")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		printHelp(stderr)
		return 2
	}

	switch {
	case showHelp:
		printHelp(stdout)
		return 0
	case showVersion:
		fmt.Fprintf(stdout, "%s %s\n", toolName, version.Version)
		return 0
	case doUpdate:
		if err := update.SelfUpdate(version.Version); err != nil {
			fmt.Fprintf(stderr, "update failed: %v\n", err)
			return 1
		}
		return 0
	}

	files := fs.Args()
	if len(files) == 0 {
		printHelp(stdout)
		return 0
	}

	w := bufio.NewWriter(stdout)
	defer w.Flush()

	exitCode := 0
	for _, file := range files {
		if err := uncomment.ProcessFile(w, file); err != nil {
			fmt.Fprintf(stderr, "%s: %v\n", file, err)
			exitCode = 1
		}
	}
	return exitCode
}

func printHelp(w io.Writer) {
	fmt.Fprintf(w, "%s %s\n", toolName, version.Version)
	fmt.Fprintf(w, "Remove comment lines from text files and print the result to stdout.\n\n")
	fmt.Fprintf(w, "Usage:\n  %s [options] <file> [file...]\n\n", toolName)
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  -h, -help       Show this help message")
	fmt.Fprintln(w, "  -v, -version    Show version information")
	fmt.Fprintln(w, "  -u, -update     Update uncmnt to the latest release")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Lines are dropped if they are blank, contain only whitespace, or start")
	fmt.Fprintln(w, "with '!', '#', or ';' (leading whitespace ignored). Multiple files are")
	fmt.Fprintln(w, "each processed and their results printed to stdout in order.")
}

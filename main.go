package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sig9org/uncmnt/internal/config"
	"github.com/sig9org/uncmnt/internal/uncomment"
	"github.com/sig9org/uncmnt/internal/update"
	"github.com/sig9org/uncmnt/internal/version"
)

const toolName = "uncmnt"

// ANSI color codes for message-level output, per uncmnt's display spec:
// normal messages are uncolored, warnings orange, errors red, debug gray.
const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[31m"
	ansiOrange = "\033[38;5;208m"
	ansiGray   = "\033[90m"
)

func colorize(color, s string) string {
	return color + s + ansiReset
}

func errorf(w io.Writer, format string, args ...any) {
	fmt.Fprintln(w, colorize(ansiRed, fmt.Sprintf(format, args...)))
}

func warnf(w io.Writer, format string, args ...any) {
	fmt.Fprintln(w, colorize(ansiOrange, fmt.Sprintf(format, args...)))
}

// newDebugLogger returns a debug-message emitter that is a no-op unless
// enabled. Active debug messages are timestamped, colored gray, and written
// to w regardless of what stream the caller would otherwise use, per spec.
func newDebugLogger(w io.Writer, enabled bool) func(format string, args ...any) {
	return func(format string, args ...any) {
		if !enabled {
			return
		}
		ts := time.Now().Format("2006-01-02T15:04:05.000Z07:00")
		msg := fmt.Sprintf(format, args...)
		fmt.Fprintln(w, colorize(ansiGray, fmt.Sprintf("[%s] DEBUG: %s", ts, msg)))
	}
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	w := bufio.NewWriter(stdout)
	defer w.Flush()

	if len(args) == 0 {
		printHelp(w)
		return 0
	}

	var showHelp, showVersion, doUpdate, debugMode bool
	var configPath string
	fs := flag.NewFlagSet(toolName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&showHelp, "h", false, "")
	fs.BoolVar(&showHelp, "help", false, "")
	fs.BoolVar(&showVersion, "v", false, "")
	fs.BoolVar(&showVersion, "version", false, "")
	fs.BoolVar(&doUpdate, "update", false, "")
	fs.BoolVar(&debugMode, "debug", false, "")
	fs.StringVar(&configPath, "c", "", "")
	fs.StringVar(&configPath, "config", "", "")

	if err := fs.Parse(args); err != nil {
		errorf(stderr, "%v", err)
		printHelp(stderr)
		return 2
	}

	dbg := newDebugLogger(w, debugMode)
	dbg("parsed flags: help=%v version=%v update=%v debug=%v", showHelp, showVersion, doUpdate, debugMode)

	switch {
	case showHelp:
		printHelp(w)
		return 0
	case showVersion:
		fmt.Fprintf(w, "%s %s\n", toolName, version.Full())
		return 0
	case doUpdate:
		dbg("starting self-update from current version %s", version.Version)
		if err := update.SelfUpdate(version.Version, dbg); err != nil {
			errorf(stderr, "update failed: %v", err)
			return 1
		}
		return 0
	}

	files := fs.Args()
	if len(files) == 0 {
		printHelp(w)
		return 0
	}

	var cfg *config.Config
	var cfgPath string
	var err error
	if configPath != "" {
		dbg("loading config from explicit -config path %q", configPath)
		cfg, err = config.LoadFrom(configPath)
		cfgPath = configPath
	} else {
		cfg, cfgPath, err = config.Load()
	}
	if err != nil {
		errorf(stderr, "%v", err)
		return 1
	}
	rules := uncomment.DefaultRules()
	switch {
	case cfg != nil && len(cfg.CommentPrefixes) > 0:
		rules = uncomment.Rules{Prefixes: cfg.CommentPrefixes}
		dbg("loaded comment rules from %s: %v", cfgPath, rules.Prefixes)
	case cfg != nil:
		dbg("config file %s defines no comment_prefixes, using default comment rules: %v", cfgPath, rules.Prefixes)
	default:
		dbg("no config file found, using default comment rules: %v", rules.Prefixes)
	}

	exitCode := 0
	for _, file := range files {
		dbg("processing file %q", file)
		if err := uncomment.ProcessFile(w, file, rules); err != nil {
			errorf(stderr, "%s: %v", file, err)
			exitCode = 1
			continue
		}
		dbg("finished file %q", file)
	}
	return exitCode
}

func printHelp(w io.Writer) {
	fmt.Fprintf(w, "%s %s\n", toolName, version.Full())
	fmt.Fprintf(w, "Remove comment lines from text files and print the result to stdout.\n\n")
	fmt.Fprintf(w, "Usage:\n  %s [options] <file> [file...]\n\n", toolName)
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  -h, -help       Show this help message")
	fmt.Fprintln(w, "  -v, -version    Show version information")
	fmt.Fprintln(w, "  -update         Update uncmnt to the latest release")
	fmt.Fprintln(w, "  -debug          Print detailed debug output to stdout")
	fmt.Fprintln(w, "  -c, -config     Force uncmnt to use this config file path")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Lines are dropped if they are blank, contain only whitespace, or start")
	fmt.Fprintln(w, "with '!', '#', or ';' (leading whitespace ignored). This can be overridden")
	fmt.Fprintln(w, "by a config.yml/config.yaml (see -c/-config, or README.md). Multiple")
	fmt.Fprintln(w, "files are each processed and their results printed to stdout in order.")
}

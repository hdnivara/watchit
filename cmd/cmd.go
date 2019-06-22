package cmd

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

// Command represents user-given arguments.
type Command struct {
	// Directories to watch.
	Dirs []string

	// Commands to execute on watched files.
	Cmds []string

	// Recursively watch directories?
	Recursive bool

	// File extenstions to monitor. Default: "md". Example: md txt.
	ext []string

	// Regex to filter files. Default: .+\.md$
	Regex string

	// Print help text.
	help bool
}

func usage(printDesc bool) {
	desc := "Watch directories for changes and run commands on changed files."
	usage := fmt.Sprintf("%s -d DIR [...] -c CMD [...] [OPTIONS]", os.Args[0])
	usageLong := `
REQUIRED
    -d, --dirs DIR [...]    Directories to watch.
    -c, --cmds CMD [...]    Commands to run on changed files.

OPTIONS
    -e, --ext EXT [...]     Watch only files with EXT extensions. Default 'md'.
    -R, --recursive         Watch directories recursively. Default: False.
    -V, --version           Print version.
    -h, --help              Print this help text.
`

	if printDesc {
		fmt.Fprintf(os.Stderr, "%s -- %s\n", os.Args[0], desc)
	}
	fmt.Fprintf(os.Stderr, "\nUSAGE: %s\n", usage)
	fmt.Fprintf(os.Stderr, "%s\n", usageLong)
}

func errExit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	usage(false)
	os.Exit(1)
}

// Parse parses and validates user-given arguments.
func Parse() *Command {
	var cmd Command

	cmdf := flag.NewFlagSet("watchit", flag.ContinueOnError)

	cmdf.StringArrayVarP(&cmd.Dirs, "dirs", "d", []string{}, "")
	cmdf.StringArrayVarP(&cmd.Cmds, "cmds", "c", []string{}, "")
	cmdf.StringArrayVarP(&cmd.ext, "ext", "e", []string{"md"}, "")
	cmdf.BoolVarP(&cmd.Recursive, "recursive", "R", false, "")
	cmdf.BoolVarP(&cmd.help, "help", "h", false, "Print help text")

	if err := cmdf.Parse(os.Args); err != nil {
		errExit(err)
	}

	if len(cmd.Dirs) == 0 || len(cmd.Cmds) == 0 {
		errExit(fmt.Errorf("-d DIR and -c CMD are required"))
	}

	if cmd.help {
		usage(true)
		os.Exit(0)
	}

	// Create regex for file extensions.
	cmd.Regex = fmt.Sprintf(".+\\.(%s)$", strings.Join(cmd.ext, "|"))

	return &cmd
}

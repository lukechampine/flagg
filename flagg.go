package flagg

import (
	"flag"
	"os"
	"unsafe"
)

// Root is the default root flag.FlagSet.
var Root = flag.CommandLine

// New returns a new flag.FlagSet. It is equivalent to flag.NewFlagSet(name,
// flag.ExitOnError).
func New(name string) *flag.FlagSet {
	return flag.NewFlagSet(name, flag.ExitOnError)
}

// SimpleUsage returns a func that writes usage to os.Stderr. If cmd has
// associated flags, the func also calls cmd.PrintDefaults.
func SimpleUsage(cmd *flag.FlagSet, usage string) func() {
	return func() {
		os.Stderr.WriteString(usage)
		numFlags := 0
		cmd.VisitAll(func(*flag.Flag) { numFlags++ })
		if numFlags > 0 {
			os.Stderr.WriteString("\nFlags:\n")
			cmd.PrintDefaults()
		}
	}
}

// A Tree is a tree of commands and subcommands.
type Tree struct {
	Cmd *flag.FlagSet
	Sub []Tree
}

// Parse parses os.Args according to the supplied Tree. It returns the
// most deeply-nested flag.FlagSet selected by the args.
func Parse(tree Tree) *flag.FlagSet {
	return parse(tree, os.Args[1:])
}

func parse(tree Tree, args []string) *flag.FlagSet {
	tree.Cmd.Parse(args)
	args = tree.Cmd.Args()
	if len(args) > 0 {
		for _, t := range tree.Sub {
			if name(t.Cmd) == args[0] {
				return parse(t, args[1:])
			}
		}
	}
	return tree.Cmd
}

// unfortunately the flag API doesn't provide a way to access the name of a
// FlagSet after it's created, so we have to resort to this hack for now.
func name(f *flag.FlagSet) string {
	// from src/flag/flag.go
	// This should be relatively safe, since the first two fields of
	// flag.FlagSet haven't changed in 6 years (!!)
	return (*struct {
		usage func()
		name  string
	})(unsafe.Pointer(f)).name
}

// Package flagg provides a simple means of constructing CLI command
// hierarchies. A hierarchy is simply a tree of flag.FlagSets; the Parse
// function is then used to determine which command is being invoked.
package flagg

import (
	"flag"
	"os"
)

// Root is the default root flag.FlagSet.
var Root = flag.CommandLine

// New returns a new flag.FlagSet. It is equivalent to flag.NewFlagSet(name,
// flag.ExitOnError), and setting f.Usage = SimpleUsage(f, usage)
func New(name, usage string) *flag.FlagSet {
	f := flag.NewFlagSet(name, flag.ExitOnError)
	f.Usage = SimpleUsage(f, usage)
	return f
}

// SimpleUsage returns a func that writes usage to cmd.Output(). If cmd has
// associated flags, the func also calls cmd.PrintDefaults.
func SimpleUsage(cmd *flag.FlagSet, usage string) func() {
	return func() {
		cmd.Output().Write([]byte(usage))
		numFlags := 0
		cmd.VisitAll(func(*flag.Flag) { numFlags++ })
		if numFlags > 0 {
			cmd.Output().Write([]byte("\nFlags:\n"))
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
			if t.Cmd.Name() == args[0] {
				return parse(t, args[1:])
			}
		}
	}
	return tree.Cmd
}

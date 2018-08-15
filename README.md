flagg
-----

[![GoDoc](https://godoc.org/lukechampine.com/flagg?status.svg)](https://godoc.org/lukechampine.com/flagg)
[![Go Report Card](http://goreportcard.com/badge/lukechampine.com/flagg)](https://goreportcard.com/report/lukechampine.com/flagg)

```
go get lukechampine.com/flagg
```

`flagg` is a tiny package that makes it easier to build CLI apps that use
subcommands. I built it because the stdlib `flag` package is too low-level,
but popular alternatives like `spf13/cobra` are full-fledged frameworks with
too many bells and whistles for my liking.

`flagg` is designed around the stdlib `*flag.FlagSet` type. Simply put, it is
a utility for constructing a hierarchy of `*flag.FlagSet`-based commands,
parsing them, and identifying which command was selected.

**NOTE: `flagg` requires Go 1.10, which added new methods to the flag.FlagSet
type.**

## Example

```go
// commands are just *flag.FlagSets
var rootCmd *flag.FlagSet = flagg.Root
rootCmd.Usage = flagg.SimpleUsage(rootCmd, "An example program")
verbose := rootCmd.Bool("v", false, "display verbose output")

// flagg.New constructs a *flag.FlagSet with the given name and usage
// description
fooCmd := flagg.New("foo", "The foo subcommand")
bars := fooCmd.Int("bars", 0, "number of bars to foo")
quuxCmd := flagg.New("quux", "The quux subcommand")

// construct the command hierarchy
tree := flagg.Tree{
	Cmd: rootCmd,
	Sub: []flagg.Tree{
		{Cmd: fooCmd},
		{Cmd: quuxCmd},
	},
}

// Parse the command hierarchy. cmd is the selected command, e.g. if
// os.Args = []string{"./example", "quux"}, then cmd == quuxCmd.
cmd := flagg.Parse(tree)

// again, cmd is just a *flag.FlagSet, so Args() returns the arguments of the
// selected command.
args := cmd.Args()

// use a switch to identify the cmd that was selected
switch cmd {
case fooCmd:
	fmt.Printf("fooing %v bars\n", *bars)
case quuxCmd:
	fmt.Println("quux!", args)
}
```

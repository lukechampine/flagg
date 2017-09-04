flagg
-----

[![GoDoc](https://godoc.org/github.com/lukechampine/flagg?status.svg)](https://godoc.org/github.com/lukechampine/flagg)
[![Go Report Card](http://goreportcard.com/badge/github.com/lukechampine/flagg)](https://goreportcard.com/report/github.com/lukechampine/flagg)

```
go get github.com/lukechampine/flagg
```

`flagg` is a tiny package that makes it easier to build a CLI app that uses
subcommands. I built it because the stdlib `flag` package is too low-level,
but popular alternatives like `spf13/cobra` are full-fledged frameworks with
too many bells and whistles for my liking.

`flagg` is designed around the stdlib `*flag.FlagSet` type. Simply put, it is
a utility for constructing a hierarchy of `*flag.FlagSet`-based commands,
parsing them, and identifying which command was selected.

## Example

```go
// commands are just *flag.FlagSets
var rootCmd *flag.FlagSet = flagg.Root
rootCmd.Usage = flagg.SimpleUsage(rootCmd, "An example program")
verbose := rootCmd.Bool("v", false, "display verbose output")

var fooCmd *flag.FlagSet = flagg.New("foo")
fooCmd.Usage = flagg.SimpleUsage(fooCmd, "The foo subcommand")
bars := fooCmd.Int("bars", 0, "number of foo'd bars")

var quuxCmd *flag.FlagSet = flagg.New("quux")
quuxCmd.Usage = flagg.SimpleUsage(quuxCmd, "The quux subcommand")

// construct command tree and parse it
// cmd is the selected command
cmd := flagg.Parse(flag.Tree{
	Cmd: rootCmd,
	Sub: []flagg.Tree{
		{Cmd: fooCmd},
		{Cmd: quuxCmd},
	},
})
// again, cmd is just a *flag.FlagSet
args := cmd.Args()

// use a switch to identify the cmd that was selected
switch cmd {
case fooCmd:
	fmt.Printf("fooing %v bars", *bars)

case quuxCmd:
	fmt.Println("quux!", args)
}
```

package flagg

import (
	"bytes"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	foo := New("foo", "")
	v := foo.Bool("v", false, "")
	bar := New("bar", "")
	os.Args = []string{"foo", "-v", "bar"}
	cmd := Parse(Tree{
		Cmd: foo,
		Sub: []Tree{
			{Cmd: bar},
		},
	})
	if cmd != bar {
		t.Fatal("expected bar, got", cmd)
	}
	if !*v {
		t.Fatal("expected -v to be set")
	}
}

func TestSimpleUsage(t *testing.T) {
	foo := New("foo", "")
	usage := "foo bar baz"
	buf := new(bytes.Buffer)
	foo.SetOutput(buf)

	// no flags
	SimpleUsage(foo, usage)()
	if out := string(buf.Next(len(usage))); out != usage {
		t.Fatalf("expected %q, got %q", usage, out)
	}

	// with flag
	foo.Bool("v", false, "")
	SimpleUsage(foo, usage)()
	exp := usage + `
Flags:
  -v	
`
	if out := string(buf.Next(len(exp))); out != exp {
		t.Fatalf("expected %q, got %q", exp, out)
	}
}

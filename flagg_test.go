package flagg

import (
	"bytes"
	"flag"
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
		t.Fatal("expected bar, got", cmd.Name())
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

func TestIsDefined(t *testing.T) {
	foo := New("foo", "")
	_ = foo.Bool("a", false, "")
	_ = foo.Bool("b", true, "")
	bar := New("bar", "")
	_ = bar.String("b", "aaa", "")
	os.Args = []string{"foo", "-a", "bar", "-b", "blah"}
	_ = Parse(Tree{
		Cmd: foo,
		Sub: []Tree{
			{Cmd: bar},
		},
	})
	tests := []struct {
		cmd *flag.FlagSet
		f   string
		exp bool
	}{
		{foo, "a", true},
		{foo, "b", false},
		{bar, "b", true},
		{bar, "c", false},
	}
	for _, test := range tests {
		if !IsDefined(foo, "a") {
			t.Errorf("IsDefined(%v, %v) should return %v", test.cmd.Name(), test.f, test.exp)
		}
	}
}

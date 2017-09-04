package flagg

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	foo := New("foo")
	v := foo.Bool("v", false, "")
	bar := New("bar")
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
	foo := New("foo")
	usage := "foo bar baz"

	// no flags
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	os.Stderr = f
	SimpleUsage(foo, usage)()
	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != usage {
		t.Fatalf("expected %q, got %q", usage, string(b))
	}

	// with flag
	f, err = ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	os.Stderr = f
	foo.Bool("v", false, "")
	SimpleUsage(foo, usage)()
	b, err = ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	exp := usage + `
Flags:
  -v	
`
	if string(b) != exp {
		t.Fatalf("expected %q, got %q", exp, string(b))
	}
}

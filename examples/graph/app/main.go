package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/closure"
	"github.com/mandelsoft/flagutils/examples/graph/graph"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/output/tableoutput"
	"github.com/mandelsoft/flagutils/sort"
)

func Error(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+msg+"\n", args...)
	os.Exit(1)
}

func main() {
	ctx := context.Background()
	opts := flagutils.DefaultOptionSet{}
	opts.Add(
		closure.NewByFactory[*graph.Element](graph.ClosureFactory),
		sort.New(),
		tableoutput.New(),
		output.New(graph.OutputsFactory),
	)

	g := New()
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	opts.AddFlags(fs)

	err := fs.Parse([]string{"c", "-s", "-value", "-c", "-o", "tree"})
	if err != nil {
		Error("%s", err)
	}

	err = flagutils.Validate(ctx, opts, nil)
	if err != nil {
		Error("%s", err)
	}

	args := fs.Args()
	out := output.From[*graph.Element](opts)
	n, err := out.GetOutput().Process(ctx, args, graph.NewSourceFactory(g))
	if err != nil {
		Error("%s", err)
	}
	fmt.Printf("processed %d nodes\n", n)
}

func New() *graph.Graph {
	a := graph.NewNode("a", "alice")
	b := graph.NewNode("b", "bob")
	c := graph.NewNode("c", "charly")
	d := graph.NewNode("d", "david")
	e := graph.NewNode("e", "eve")

	a.AddChild(d)
	a.AddChild(e)
	b.AddChild(d)
	c.AddChild(b)
	c.AddChild(a)
	c.AddChild(e)
	e.AddChild(d)

	b.AddChild(c)

	g := graph.NewGraph()
	g.AddRoot(c)

	return g
}

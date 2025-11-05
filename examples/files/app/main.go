package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/closure"
	"github.com/mandelsoft/flagutils/examples/files/files"
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
		files.New(),
		closure.New[*files.Element](files.Closure),
		sort.New(),
		tableoutput.New(),
		output.New(files.OutputsFactory),
	)

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	opts.AddFlags(fs)

	err := fs.Parse([]string{"-c", "-s", "name", "output", "examples", "-o", "tree"})
	// err := fs.Parse([]string{"-c", "-s", "name", "output", "examples", "-o", "wide"})
	//err := fs.Parse([]string{"-c", "output", "examples", "-o", "YAML"})
	if err != nil {
		Error("%s", err)
	}

	err = flagutils.Validate(ctx, opts, nil)
	if err != nil {
		Error("%s", err)
	}

	args := fs.Args()
	out := output.From[*files.Element](opts)
	n, err := out.GetOutput().Process(ctx, args, files.NewSourceFactory(opts))
	if err != nil {
		Error("%s", err)
	}
	fmt.Printf("processed %d files\n", n)
}

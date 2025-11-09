package manifest

import (
	"context"
	output "github.com/mandelsoft/flagutils/output/internal"
	"github.com/mandelsoft/flagutils/utils/out"
	"github.com/mandelsoft/goutils/iterutils"
	"github.com/mandelsoft/streaming"
	"iter"
)

// Factory is a ProcessorFactory and Processor in one type
// because no state is required.
type Factory struct {
	formatter Formatter
}

var _ streaming.ProcessorFactory[output.ElementSpecs, output.Result, Manifest] = (*Factory)(nil)

func (o *Factory) Processor(output.ElementSpecs) (streaming.Processor[output.Result, Manifest], error) {
	return o.Process, nil
}

var (
	_ streaming.Processor[output.Result, Manifest] = (*Factory)(nil).Process
)

func (p *Factory) Process(ctx context.Context, i iter.Seq[Manifest]) (int, error) {
	d := iterutils.Get(i)

	if len(d) == 0 {
		out.Print(ctx, "no elements found\n")
		return 0, nil
	}
	p.formatter.Format(ctx, d)
	return len(d), nil
}

type wrapper struct {
	e any
}

func (w *wrapper) AsManifest() any {
	return w.e
}

func mapToManifest[I any](in I) Manifest {
	if m, ok := any(in).(Manifest); ok {
		return m
	} else {
		return &wrapper{in}
	}
}

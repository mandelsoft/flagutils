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
type Factory[I any] struct {
	factory *OutputFactory[I]
}

var _ streaming.ProcessorFactory[output.ElementSpecs, output.Result, any] = (*Factory[any])(nil)

func (o *Factory[I]) Processor(output.ElementSpecs) (streaming.Processor[output.Result, I], error) {
	return o.Process, nil
}

var (
	_ streaming.Processor[output.Result, any] = (*Factory[any])(nil).Process
)

type wrapper struct {
	e any
}

func (w *wrapper) AsManifest() any {
	return w.e
}

func (p *Factory[I]) Process(ctx context.Context, i iter.Seq[I]) (int, error) {
	d := iterutils.Get(i)

	if len(d) == 0 {
		out.Print(ctx, "no elements found\n")
		return 0, nil
	}
	r := make([]Manifest, len(d))
	for i, e := range d {
		var o any = e
		if m, ok := o.(Manifest); ok {
			r[i] = m
		} else {
			r[i] = &wrapper{e}
		}
	}
	p.factory.formatter.Format(ctx, r)
	return len(d), nil
}

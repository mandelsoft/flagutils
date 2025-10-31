package tableoutput

import (
	"context"
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/sort"
	"github.com/mandelsoft/streaming/chain"
	"slices"
)

func NewOutputFactory[I any](mapper chain.Mapper[I, []string], headers ...string) output.OutputFactory[I] {
	return &OutputFactory[I]{mapper, slices.Clone(headers)}
}

type OutputFactory[I any] struct {
	mapper  chain.Mapper[I, []string]
	headers []string
}

var _ output.OutputFactory[int] = (*OutputFactory[int])(nil)

func (o *OutputFactory[I]) GetFieldNames() []string {
	return o.headers
}

func (o *OutputFactory[I]) Create(ctx context.Context, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (output.Output[I], error) {
	c := chain.Mapped[I, []string](o.mapper)
	s := sort.From(opts)
	if s != nil {
		c = sort.AddSortChain(c, s)
	}
	return output.NewOutput[I, []string](c, &Factory{slices.Clone(o.headers), From(opts)}), nil
}

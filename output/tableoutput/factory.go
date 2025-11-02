package tableoutput

import (
	"context"
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/closure"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/sort"
	"github.com/mandelsoft/streaming/chain"
	"slices"
	"strings"
)

type FieldProvider = output.FieldProvider

func NewOutputFactory[I any](mapper chain.Mapper[I, FieldProvider], headers ...string) *OutputFactory[I, I] {
	return &OutputFactory[I, I]{chain.New[I](), mapper, slices.Clone(headers)}
}

func NewExtendedOutputFactory[I, O any](chain chain.Chain[I, O], mapper chain.Mapper[O, FieldProvider], headers ...string) *OutputFactory[I, O] {
	return &OutputFactory[I, O]{chain, mapper, slices.Clone(headers)}
}

type OutputFactory[I, O any] struct {
	chain   chain.Chain[I, O]
	mapper  chain.Mapper[O, FieldProvider]
	headers []string
}

var _ output.OutputFactory[int] = (*OutputFactory[int, int])(nil)

func (o *OutputFactory[I, O]) GetMapper() chain.Mapper[O, FieldProvider] {
	return o.mapper
}

func (o *OutputFactory[I, O]) GetHeaders() []string {
	return slices.Clone(o.headers)
}

func (o *OutputFactory[I, O]) GetFieldNames() []string {
	fields := slices.Clone(o.headers)
	for i := range fields {
		if strings.HasPrefix(fields[i], "-") {
			fields[i] = fields[i][1:]
		}
	}
	return fields
}

func (o *OutputFactory[I, O]) Create(ctx context.Context, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (output.Output[I], error) {
	c := chain.New[I]()

	e := closure.From[I](opts)
	if e != nil && e.GetExploder() != nil {
		c = chain.AddExplode(c, e.GetExploder())
	}

	co := chain.AddChain(c, o.chain)
	mapped := chain.AddMap[FieldProvider](co, o.mapper)
	s := sort.From(opts)
	if s != nil {
		mapped = sort.AddSortChain(mapped, s)
	}
	return output.NewOutput[I, FieldProvider](mapped, &Factory{slices.Clone(o.headers), From(opts)}), nil
}

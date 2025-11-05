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

func NewOutputFactory[I any](mapper chain.Mapper[I, FieldProvider], headers ...string) *OutputFactory[I, FieldProvider] {
	return &OutputFactory[I, FieldProvider]{mapper, chain.New[FieldProvider](), slices.Clone(headers)}
}

func NewExtendedOutputFactory[I any, F FieldProvider](mapper chain.Mapper[I, F], chain chain.Chain[F, FieldProvider], headers ...string) *OutputFactory[I, F] {
	return &OutputFactory[I, F]{mapper, chain, slices.Clone(headers)}
}

type OutputFactory[I any, F FieldProvider] struct {
	mapper  chain.Mapper[I, F]
	chain   chain.Chain[F, FieldProvider]
	headers []string
}

var _ output.OutputFactory[int] = (*OutputFactory[int, FieldProvider])(nil)

func (o *OutputFactory[I, F]) GetMapper() chain.Mapper[I, F] {
	return o.mapper
}

func (o *OutputFactory[I, F]) GetHeaders() []string {
	return slices.Clone(o.headers)
}

func (o *OutputFactory[I, F]) GetFieldNames(stage string) []string {
	fields := slices.Clone(o.headers)
	for i := range fields {
		if strings.HasPrefix(fields[i], "-") {
			fields[i] = fields[i][1:]
		}
	}
	return fields
}

func (o *OutputFactory[I, F]) Create(ctx context.Context, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (output.Output[I], error) {
	c := chain.New[I]()

	e := closure.From[I](opts)
	if e != nil {
		f := e.GetExploderFactory(opts)
		if f != nil {
			c = chain.AddExplodeByFactory(c, f)
		}
	}
	mapped := chain.AddMap[F](c, o.mapper)
	s := sort.From(opts)
	if s != nil {
		mapped = sort.AddSortChain[I, F](mapped, s)
	}

	co := chain.AddChain(mapped, o.chain)
	return output.NewOutput[I, FieldProvider](co, &Factory{slices.Clone(o.headers), From(opts)}), nil
}

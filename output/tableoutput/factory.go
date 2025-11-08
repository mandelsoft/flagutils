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
type ExtendableFieldProvider = output.ExtendableFieldProvider

func Cast[O, I any](o I) O {
	return any(o).(O)
}

func NewOutputFactory[I any, F FieldProvider](mapper chain.Mapper[I, F], headers ...string) *OutputFactory[I, F] {
	return &OutputFactory[I, F]{mapper: mapper, chain: chain.Mapped[F, FieldProvider](Cast[FieldProvider, F]), headers: slices.Clone(headers)}
}

func NewOutputFactoryByProvider[I any, F FieldProvider](provider output.MappingProvider[I, F]) *OutputFactory[I, F] {
	return &OutputFactory[I, F]{provider: provider, chain: chain.Mapped[F, FieldProvider](Cast[FieldProvider, F])}
}

func NewExtendedOutputFactory[I any, F FieldProvider](mapper chain.Mapper[I, F], chain chain.Chain[F, FieldProvider], headers ...string) *OutputFactory[I, F] {
	return &OutputFactory[I, F]{mapper: mapper, chain: chain, headers: slices.Clone(headers)}
}

type OutputFactory[I any, F FieldProvider] struct {
	provider output.MappingProvider[I, F]
	mapper   chain.Mapper[I, F]
	chain    chain.Chain[F, FieldProvider]
	headers  []string
}

var _ output.OutputFactory[int] = (*OutputFactory[int, FieldProvider])(nil)

func (o *OutputFactory[I, F]) GetMapper() chain.Mapper[I, F] {
	return o.mapper
}

func (o *OutputFactory[I, F]) GetProvider() output.MappingProvider[I, F] {
	return o.provider
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
			c = chain.AddExplodeByFactory[I](c, f)
		}
	}
	mapper := o.mapper
	if mapper == nil {
		var err error
		mapper, o.headers, err = o.provider.GetMapping(opts)
		if err != nil {
			return nil, err
		}
	}
	mapped := chain.AddMap[F](c, mapper)
	s := sort.From(opts)
	if s != nil {
		mapped = sort.AddSortChain[I, F](mapped, s)
	}

	co := chain.AddChain(mapped, o.chain)
	return output.NewOutput[I, FieldProvider](co, &Factory[FieldProvider]{slices.Clone(o.headers), From(opts)}), nil
}

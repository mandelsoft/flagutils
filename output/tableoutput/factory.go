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
func (o *OutputFactory[I, F]) getMapper(opts flagutils.OptionSetProvider) (chain.Mapper[I, F], error) {
	mapper := o.mapper
	if mapper == nil {
		var err error
		mapper, o.headers, err = o.provider.GetMapping(opts)
		if err != nil {
			return nil, err
		}
	}
	return mapper, nil
}

func (o *OutputFactory[I, F]) Create(ctx context.Context, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (output.Output[I], error) {
	mapper, err := o.getMapper(opts)
	if err != nil {
		return nil, err
	}

	// compose chain: exploder -> mapper -> sort -> custom chain
	c := closure.AddExplodeChain(opts, chain.New[I]())
	mapped := sort.AddSortChain[I, F](opts, chain.AddMap[F](c, mapper))
	co := chain.AddChain(mapped, o.chain)
	return output.NewOutput[I, FieldProvider](co, &Factory[FieldProvider]{slices.Clone(o.headers), From(opts)}), nil
}

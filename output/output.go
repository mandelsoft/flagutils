package output

import (
	"context"
	"github.com/mandelsoft/flagutils"
	"slices"
	"strings"

	"github.com/mandelsoft/streaming"
	"github.com/mandelsoft/streaming/chain"
)

type Result = int

////////////////////////////////////////////////////////////////////////////////

type Output[I any] interface {
	FieldNameProvider
	Process(ctx context.Context, specs ElementSpecs, src streaming.SourceFactory[ElementSpecs, I]) (Result, error)
}

type DefaultOutput[I, O any] struct {
	fieldNames []string
	chain      chain.Chain[I, O]
	processor  streaming.ProcessorFactory[ElementSpecs, Result, O]
}

var _ Output[string] = (*DefaultOutput[string, []string])(nil)

func NewOutput[I, O any](chain chain.Chain[I, O], processor streaming.ProcessorFactory[ElementSpecs, Result, O]) *DefaultOutput[I, O] {
	return &DefaultOutput[I, O]{nil, chain, processor}
}

func (o *DefaultOutput[I, O]) GetFieldNames() []string {
	fields := slices.Clone(o.fieldNames)
	for i := range fields {
		if strings.HasPrefix(fields[i], "-") {
			fields[i] = fields[i][1:]
		}
	}
	return fields
}

func (o *DefaultOutput[I, O]) GetChain() chain.Chain[I, O] {
	return o.chain
}

func (o *DefaultOutput[I, O]) GetProcessor() streaming.ProcessorFactory[ElementSpecs, Result, O] {
	return o.processor
}

func (o *DefaultOutput[I, O]) Process(ctx context.Context, specs ElementSpecs, src streaming.SourceFactory[ElementSpecs, I]) (Result, error) {
	s, err := src.Elements(specs)
	if err != nil {
		return 0, err
	}
	return streaming.NewSink[ElementSpecs, Result, I, O](o.chain, o.processor).Execute(ctx, specs, s)
}

////////////////////////////////////////////////////////////////////////////////

type OutputFactory[I any] interface {
	FieldNameProvider
	Create(ctx context.Context, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (Output[I], error)
}

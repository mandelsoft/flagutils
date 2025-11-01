package internal

import (
	"context"
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/streaming"
)

type FieldNameProvider interface {
	GetFieldNames() []string
}

type Result = int

type ElementSpecs interface{}

////////////////////////////////////////////////////////////////////////////////

type OutputFactory[I any] interface {
	FieldNameProvider
	Create(ctx context.Context, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (Output[I], error)
}

type Output[I any] interface {
	FieldNameProvider
	Process(ctx context.Context, specs ElementSpecs, src streaming.SourceFactory[ElementSpecs, I]) (Result, error)
}

////////////////////////////////////////////////////////////////////////////////

type OutputsFactory[I any] interface {
	GetModes() []string
	Add(mode string, out OutputFactory[I]) OutputsFactory[I]
	AddManifestOutputs() OutputsFactory[I]

	GetFieldNames(mode string) []string
	CreateOutput(ctx context.Context, mode string, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (Output[I], error)
}

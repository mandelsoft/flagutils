package internal

import (
	"context"
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/streaming"
)

type Fields []string

var _ FieldProvider = Fields{}

func (f Fields) GetFields() []string {
	return f
}

type FieldProvider interface {
	GetFields() []string
}

type FieldNameProvider interface {
	GetFieldNames(stage string) []string
}

type Result = int

type ElementSpecs interface{}

////////////////////////////////////////////////////////////////////////////////

type OutputFactory[I any] interface {
	FieldNameProvider
	Create(ctx context.Context, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (Output[I], error)
}

type Output[I any] interface {
	Process(ctx context.Context, specs ElementSpecs, src streaming.SourceFactory[ElementSpecs, I]) (Result, error)
}

////////////////////////////////////////////////////////////////////////////////

type OutputsFactory[I any] interface {
	GetModes() []string
	Add(mode string, out OutputFactory[I]) OutputsFactory[I]
	AddManifestOutputs() OutputsFactory[I]

	GetFieldNames(mode, stage string) []string
	CreateOutput(ctx context.Context, mode string, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (Output[I], error)
}

package internal

import (
	"context"
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/streaming"
	"github.com/mandelsoft/streaming/chain"
	"slices"
)

type Fields []string

var _ ExtendedFieldProvider = (*Fields)(nil)

func (f Fields) GetFields() []string {
	return f
}

func (f *Fields) InsertFields(i int, s ...string) {
	*f = slices.Insert(*f, i, s...)
}

type FieldProvider interface {
	GetFields() []string
}

type ExtendedFieldProvider interface {
	FieldProvider
	InsertFields(int, ...string)
}

type FieldNameProvider interface {
	GetFieldNames(stage string) []string
}

type MappingProvider[I any, F FieldProvider] interface {
	// GetMapping provides a mapper of an element to fields
	// and the appropriate header fields.
	GetMapping(opts flagutils.OptionSetProvider) (chain.Mapper[I, F], []string, error)
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

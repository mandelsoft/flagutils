package output

import (
	"github.com/mandelsoft/flagutils/output/internal"
)

type Fields = internal.Fields

type FieldNameProvider = internal.FieldNameProvider
type FieldProvider = internal.FieldProvider
type ExtendableFieldProvider = internal.ExtendedFieldProvider
type ElementSpecs = internal.ElementSpecs
type Result = internal.Result

type MappingProvider[I any, F FieldProvider] = internal.MappingProvider[I, F]

////////////////////////////////////////////////////////////////////////////////

type OutputFactory[I any] = internal.OutputFactory[I]
type Output[I any] = internal.Output[I]

type OutputsFactory[I any] = internal.OutputsFactory[I]

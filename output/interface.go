package output

import (
	"github.com/mandelsoft/flagutils/output/internal"
)

type FieldNameProvider = internal.FieldNameProvider
type ElementSpecs = internal.ElementSpecs
type Result = internal.Result

////////////////////////////////////////////////////////////////////////////////

type OutputFactory[I any] = internal.OutputFactory[I]
type Output[I any] = internal.Output[I]

type OutputsFactory[I any] = internal.OutputsFactory[I]

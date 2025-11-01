package output

import (
	"github.com/mandelsoft/flagutils/output/internal"
	"github.com/mandelsoft/streaming"
	"github.com/mandelsoft/streaming/chain"
)

func NewOutput[I, O any](chain chain.Chain[I, O], processor streaming.ProcessorFactory[ElementSpecs, Result, O]) *internal.DefaultOutput[I, O] {
	return internal.NewOutput(chain, processor)
}

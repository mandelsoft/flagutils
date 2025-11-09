package parallel

import (
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/streaming/chain"
)

// AddParallelChain evaluates the parallel Options (if present) to decide
// whether chain a should be executed parallel or just added
// to chain c.
// This helper can be used, for example, by output implementations
// to organize their processing chains.
func AddParallelChain[N, I, O any](opts flagutils.OptionSetProvider, c chain.Chain[I, O], a chain.Chain[O, N]) chain.Chain[I, N] {
	o := From(opts)
	if o != nil {
		p := o.GetPool()
		if p != nil {
			return chain.AddParallel[N](c, a, p)
		}
	}
	return chain.AddChain[N](c, a)
}

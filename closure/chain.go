package closure

import (
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/streaming/chain"
)

// AddExplodeChain evaluates the closure Options (if present) to decide
// whether an explode step should be added to chain c.
// This helper can be used, for example, by output implementations
// to organize their processing chains.
func AddExplodeChain[I, O any](opts flagutils.OptionSetProvider, c chain.Chain[I, O]) chain.Chain[I, O] {
	o := From[O](opts)
	if o == nil || !o.Value() {
		return c
	}
	f := o.GetExploderFactory(opts)
	if f == nil {
		return c
	}
	return chain.AddExplodeByFactory[O](c, f)
}

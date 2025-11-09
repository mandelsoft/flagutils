package sort

import (
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/streaming/chain"
)

// AddSortChain evaluates the sort Options (if present) to decide
// whether a sort step should be added to chain c.
// This helper can be used, for example, by output implementations
// to organize their processing chains.
func AddSortChain[I any, F output.FieldProvider](opts flagutils.OptionSetProvider, c chain.Chain[I, F]) chain.Chain[I, F] {
	o := From(opts)
	if o == nil || len(o.Value()) == 0 {
		return c
	}
	return chain.AddSort(c, func(a, b F) int { return o.Compare(a, b) })
}

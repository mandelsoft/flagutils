package closure

import (
	"context"
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/streaming/chain"
	"github.com/spf13/pflag"
)

type Options[I any] struct {
	closure  bool
	exploder chain.Exploder[I, I]
}

func From[I any](opts flagutils.OptionSetProvider) *Options[I] {
	return flagutils.GetFrom[*Options[I]](opts)
}

var (
	_ flagutils.Options = (*Options[int])(nil)
)

func New[I any](exploder chain.Exploder[I, I]) *Options[I] {
	return &Options[I]{exploder: exploder}
}

func (o *Options[I]) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&o.closure, "closure", "c", false, "calculate closure")
}

func (o *Options[I]) GetExploder() chain.Exploder[I, I] {
	if !o.closure {
		return nil
	}
	return o.exploder
}

func AddExploderChain[I, O any](c chain.Chain[I, O], opts *Options[O]) chain.Chain[I, O] {
	return chain.AddConditional(c,
		func(context.Context) bool { return opts != nil && opts.closure },
		chain.Exploded(opts.exploder),
	)
}

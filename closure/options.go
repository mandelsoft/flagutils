package closure

import (
	"context"
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/streaming/chain"
	"github.com/spf13/pflag"
)

type Options[I any] struct {
	flagutils.OptionBase[*Options[I]]
	closure  bool
	exploder ExploderFactory[I]
}

func From[I any](opts flagutils.OptionSetProvider) *Options[I] {
	o := flagutils.GetFrom[*Options[I]](opts)
	o.OptionBase = flagutils.NewBase(o)
	return o
}

var (
	_ flagutils.Options = (*Options[int])(nil)
)

type ExploderFactory[I any] func(opts flagutils.OptionSetProvider) chain.ExploderFactory[I, I]

func ExploderFactoryFor[I any](e chain.Exploder[I, I]) ExploderFactory[I] {
	return func(opts flagutils.OptionSetProvider) chain.ExploderFactory[I, I] {
		return chain.ExploderFactoryFor[I, I](e)
	}
}

func New[I any](exploder chain.Exploder[I, I]) *Options[I] {
	return &Options[I]{exploder: ExploderFactoryFor[I](exploder)}
}

func NewByFactory[I any](exploder ExploderFactory[I]) *Options[I] {
	return &Options[I]{exploder: exploder}
}

func (o *Options[I]) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&o.closure, "closure", "c", false, "calculate closure")
}

func (o *Options[I]) GetExploderFactory(opts flagutils.OptionSetProvider) chain.ExploderFactory[I, I] {
	if !o.closure || o.exploder == nil {
		return nil
	}
	return o.exploder(opts)
}

func AddExploderChain[I, O any](c chain.Chain[I, O], opts flagutils.OptionSetProvider) chain.Chain[I, O] {
	o := From[O](opts)
	return chain.AddConditional(c,
		func(context.Context) bool { return o != nil && o.closure },
		chain.ExplodedByFactory(o.GetExploderFactory(opts)),
	)
}

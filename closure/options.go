package closure

import (
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/streaming/chain"
)

type Options[I any] struct {
	flagutils.SimpleOption[bool, *Options[I]]
	exploder ExploderFactory[I]
}

func From[I any](opts flagutils.OptionSetProvider) *Options[I] {
	o := flagutils.GetFrom[*Options[I]](opts)
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
	return NewByFactory[I](ExploderFactoryFor[I](exploder))
}

func NewByFactory[I any](exploder ExploderFactory[I]) *Options[I] {
	o := &Options[I]{exploder: exploder}
	o.SimpleOption = flagutils.NewSimpleOption(o, false, "closure", "c", "calculate closure")
	return o
}

func (o *Options[I]) GetExploderFactory(opts flagutils.OptionSetProvider) chain.ExploderFactory[I, I] {
	if !o.Value() || o.exploder == nil {
		return nil
	}
	return o.exploder(opts)
}

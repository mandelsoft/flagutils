package parallel

import (
	"context"
	"fmt"

	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/streaming/processing"
	"github.com/mandelsoft/streaming/simplepool"
)

type Options struct {
	flagutils.SimpleOption[int, *Options]
	poolprovider PoolProvider

	pool processing.Processing
}

func From(opts flagutils.OptionSetProvider) *Options {
	return flagutils.GetFrom[*Options](opts)
}

var (
	_ flagutils.Options     = (*Options)(nil)
	_ flagutils.Validatable = (*Options)(nil)
	_ flagutils.Finalizable = (*Options)(nil)
)

type PoolProvider func(ctx context.Context, n int) processing.Processing

func New(n ...int) *Options {
	o := &Options{}
	o.SimpleOption = flagutils.NewSimpleOption[int](o, general.Optional(n...), "parallel", "p", "degree of parallelism")
	return o
}

func (o *Options) WithPoolProvider(p PoolProvider) *Options {
	o.poolprovider = p
	return o
}

func (o *Options) GetPool() processing.Processing {
	return o.pool
}

func (o *Options) Validate(ctx context.Context, opts flagutils.OptionSet, v flagutils.ValidationSet) error {
	n := o.Value()
	if n < 0 {
		return fmt.Errorf("invalid degree of parallelism: %d", n)
	}
	if o.pool == nil {
		if o.poolprovider != nil {
			o.pool = o.poolprovider(ctx, n)
		} else {
			o.pool = simplepool.New(ctx, n)
		}
	}
	return nil
}

func (o *Options) Finalize(ctx context.Context, opts flagutils.OptionSet, v flagutils.FinalizationSet) error {
	var err error
	if o.pool != nil {
		err = o.pool.Close()
		o.pool = nil
	}
	return err
}

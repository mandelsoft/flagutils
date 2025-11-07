package parallel

import (
	"context"
	"fmt"

	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/streaming/processing"
	"github.com/mandelsoft/streaming/simplepool"
	"github.com/spf13/pflag"
)

type Options struct {
	flagutils.OptionBase[*Options]
	poolprovider PoolProvider

	n int

	pool processing.Processing
}

func From(opts flagutils.OptionSetProvider) *Options {
	return flagutils.GetFrom[*Options](opts)
}

var (
	_ flagutils.Options     = (*Options)(nil)
	_ flagutils.Validation  = (*Options)(nil)
	_ flagutils.Finalizable = (*Options)(nil)
)

type PoolProvider func(ctx context.Context, n int) processing.Processing

func New(n ...int) *Options {
	o := &Options{n: general.Optional(n...)}
	o.OptionBase = flagutils.NewBase(o)
	return o
}

func (o *Options) WithPoolProvider(p PoolProvider) *Options {
	o.poolprovider = p
	return o
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.IntVarP(&o.n, o.Long("parallel"), o.Short("p"), o.n, o.Desc("degree of parallelism"))
}

func (o *Options) GetPool() processing.Processing {
	return o.pool
}

func (o *Options) Validate(ctx context.Context, opts flagutils.OptionSet, v flagutils.ValidationSet) error {
	if o.n < 0 {
		return fmt.Errorf("invalid degree of parallelism: %d", o.n)
	}
	if o.pool == nil {
		if o.poolprovider != nil {
			o.pool = o.poolprovider(ctx, o.n)
		} else {
			o.pool = simplepool.New(ctx, o.n)
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

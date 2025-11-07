package output

import (
	"context"
	"fmt"
	"strings"

	"github.com/mandelsoft/flagutils"
	"github.com/spf13/pflag"
)

func From[I any](opts flagutils.OptionSetProvider) *Options[I] {
	return flagutils.GetFrom[*Options[I]](opts)
}

type Options[I any] struct {
	flagutils.OptionBase[*Options[I]]
	// config
	factory OutputsFactory[I]

	// in
	mode string

	// out
	output Output[I]
}

var (
	_ flagutils.Options    = (*Options[int])(nil)
	_ FieldNameProvider    = (*Options[int])(nil)
	_ flagutils.Validation = (*Options[int])(nil)
)

func New[I any](out OutputsFactory[I]) *Options[I] {
	o := &Options[I]{factory: out}
	o.OptionBase = flagutils.NewBase(o)
	return o
}

func (o *Options[I]) AddFlags(fs *pflag.FlagSet) {
	keys := o.factory.GetModes()
	if len(keys) > 0 {
		fs.StringVarP(&o.mode, o.Long("mode"), o.Short("o"), "", fmt.Sprintf(o.Desc("output mode (%s)"), strings.Join(keys, ", ")))
	}
}

func (o *Options[I]) GetMode() string {
	return o.mode
}

func (o *Options[I]) GetOutputs() OutputsFactory[I] {
	return o.factory
}

func (o *Options[I]) GetOutput() Output[I] {
	return o.output
}

func (o *Options[I]) GetFieldNames(stage string) []string {
	return o.factory.GetFieldNames(o.mode, stage)
}

func (o *Options[I]) Validate(ctx context.Context, opts flagutils.OptionSet, v flagutils.ValidationSet) error {
	of, err := o.factory.CreateOutput(ctx, o.mode, opts, v)
	if err != nil {
		return err
	}
	o.output = of
	return nil
}

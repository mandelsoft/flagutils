package output

import (
	"context"
	"fmt"
	"strings"

	"github.com/mandelsoft/flagutils"
)

func From[I any](opts flagutils.OptionSetProvider) *Options[I] {
	return flagutils.GetFrom[*Options[I]](opts)
}

type Options[I any] struct {
	flagutils.SimpleOption[string, *Options[I]]
	// config
	factory OutputsFactory[I]

	// out
	output Output[I]
}

var (
	_ flagutils.Options     = (*Options[int])(nil)
	_ FieldNameProvider     = (*Options[int])(nil)
	_ flagutils.Validatable = (*Options[int])(nil)
)

func New[I any](out OutputsFactory[I]) *Options[I] {
	o := &Options[I]{factory: out}
	o.SimpleOption = flagutils.NewSimpleOption[string](o, "", "mode", "o", o.description("output mode (%s)"))
	return o
}

func (o *Options[I]) WithDescription(s string) *Options[I] {
	return o.SimpleOption.WithDescription(o.description(s))
}

func (o *Options[I]) description(msg string) string {
	keys := o.factory.GetModes()
	return fmt.Sprintf(msg, strings.Join(keys, ", "))
}

func (o *Options[I]) GetMode() string {
	return o.Value()
}

func (o *Options[I]) GetOutputs() OutputsFactory[I] {
	return o.factory
}

func (o *Options[I]) GetOutput() Output[I] {
	return o.output
}

func (o *Options[I]) GetFieldNames(stage string) []string {
	return o.factory.GetFieldNames(o.Value(), stage)
}

func (o *Options[I]) Validate(ctx context.Context, opts flagutils.OptionSet, v flagutils.ValidationSet) error {
	of, err := o.factory.CreateOutput(ctx, o.Value(), opts, v)
	if err != nil {
		return err
	}
	o.output = of
	return nil
}

package sort

import (
	"context"
	"fmt"
	"github.com/mandelsoft/goutils/general"
	"slices"
	"sort"
	"strings"

	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/streaming/chain"
	"github.com/spf13/pflag"
)

type Options struct {
	sortFields  []string
	fields      []string
	comparators map[string]general.CompareFunc[string]
}

func From(opts flagutils.OptionSetProvider) *Options {
	return flagutils.GetFrom[*Options](opts)
}

var (
	_ flagutils.Options    = (*Options)(nil)
	_ flagutils.Validation = (*Options)(nil)
)

func New() *Options {
	return &Options{comparators: make(map[string]general.CompareFunc[string])}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVarP(&o.sortFields, "sort", "s", []string{}, "sort fields")
}

func (o *Options) AddComparator(name string, cmp general.CompareFunc[string]) *Options {
	o.comparators[name] = cmp
	return o
}

func (o *Options) GetComparator(name string) general.CompareFunc[string] {
	return o.comparators[name]
}

func (o *Options) Validate(ctx context.Context, opts flagutils.OptionSet, v flagutils.ValidationSet) error {
	if len(o.sortFields) == 0 {
		return nil
	}
	for i, v := range o.sortFields {
		o.sortFields[i] = strings.ToLower(v)
	}

	fields, err := flagutils.ValidatedOptions[output.FieldNameProvider](ctx, opts, v)
	if err != nil {
		return err
	}

	if fields == nil {
		return fmt.Errorf("invalid sort fields: %v", o.sortFields)
	}
	o.fields = fields.GetFieldNames()
	if o.fields == nil {
		return fmt.Errorf("invalid sort fields: %v", o.sortFields)
	}
	for i, v := range o.fields {
		o.fields[i] = strings.ToLower(v)
	}
	var wrong []string
	var names []string = o.fields
	for _, o := range o.sortFields {
		if !slices.Contains(names, o) {
			wrong = append(wrong, o)
		}
	}
	if len(wrong) != 0 {
		sort.Strings(wrong)
		return fmt.Errorf("invalid sort fields: %v", wrong)
	}
	return nil
}

func AddSortChain[I any, F output.FieldProvider](c chain.Chain[I, F], opts *Options) chain.Chain[I, F] {
	return chain.AddConditional(c,
		func(context.Context) bool { return opts != nil && len(opts.sortFields) != 0 },
		chain.Sorted[F](func(a, b F) int { return opts.Compare(a, b) }),
	)
}

func (o *Options) Compare(af, bf output.FieldProvider) int {
	a := af.GetFields()
	b := bf.GetFields()
	for i := range o.sortFields {
		f := o.sortFields[len(o.sortFields)-i-1]
		i := slices.Index(o.fields, f)
		if i >= 0 {
			cmp := o.comparators[f]
			if cmp == nil {
				cmp = strings.Compare
			}
			if c := cmp(a[i], b[i]); c != 0 {
				return c
			}
		}
	}
	return 0
}

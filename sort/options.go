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
)

const FIELD_MODE_SORT = "<sort>"

type Options struct {
	flagutils.SimpleOption[[]string, *Options]

	fieldInfos  []*fieldInfo
	comparators map[string]general.CompareFunc[string]
}

func From(opts flagutils.OptionSetProvider) *Options {
	return flagutils.GetFrom[*Options](opts)
}

var (
	_ flagutils.Options     = (*Options)(nil)
	_ flagutils.Validatable = (*Options)(nil)
)

func New() *Options {
	o := &Options{comparators: make(map[string]general.CompareFunc[string])}
	o.SimpleOption = flagutils.NewSimpleOption[[]string](o, nil, "sort", "s", "sort fields")
	return o
}

func (o *Options) WithComparator(name string, cmp general.CompareFunc[string]) *Options {
	o.comparators[strings.ToLower(name)] = cmp
	return o
}

func (o *Options) GetComparator(name string) general.CompareFunc[string] {
	return o.comparators[name]
}

type fieldInfo struct {
	order int
	index int
	cmp   general.CompareFunc[string]
}

func (o *Options) Validate(ctx context.Context, opts flagutils.OptionSet, v flagutils.ValidationSet) error {
	sortFields := o.Value()
	if len(sortFields) == 0 {
		return nil
	}
	for i, v := range sortFields {
		sortFields[i] = strings.ToLower(v)
	}

	fields, err := flagutils.ValidatedOptions[output.FieldNameProvider](ctx, opts, v)
	if err != nil {
		return err
	}

	if fields == nil {
		return fmt.Errorf("invalid sort fields: %v", sortFields)
	}
	names := fields.GetFieldNames(FIELD_MODE_SORT)
	if names == nil {
		return fmt.Errorf("invalid sort fields: %v", sortFields)
	}
	for i, n := range names {
		names[i] = strings.ToLower(n)
	}

	var wrong []string
	for _, v := range sortFields {
		order := 1
		if strings.HasPrefix(v, "-") {
			order = -1
			v = v[1:]
		}
		idx := slices.Index(names, strings.ToLower(v))
		if idx < 0 {
			wrong = append(wrong, v)
		}
		cmp := o.comparators[v]
		if cmp == nil {
			cmp = strings.Compare
		}
		info := &fieldInfo{order: order, index: idx, cmp: cmp}
		o.fieldInfos = append(o.fieldInfos, info)
	}

	if len(wrong) != 0 {
		sort.Strings(wrong)
		return fmt.Errorf("invalid sort fields: %v", wrong)
	}
	slices.Reverse(o.fieldInfos)
	return nil
}

func (o *Options) Compare(af, bf output.FieldProvider) int {
	a := af.GetFields()
	b := bf.GetFields()
	for _, i := range o.fieldInfos {
		if c := i.cmp(a[i.index], b[i.index]); c != 0 {
			return c * i.order
		}
	}
	return 0
}

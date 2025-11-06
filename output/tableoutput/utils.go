package tableoutput

import (
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/closure"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/utils/history"
	"github.com/mandelsoft/streaming/chain"
	"slices"
)

////////////////////////////////////////////////////////////////////////////////

type FieldExtender[I, F any] interface {
	Extend(I, F) F
}

type FieldExtenderFunc[I, F any] func(o I, in F) F

func (f FieldExtenderFunc[I, F]) Extend(o I, in F) F {
	return f(o, in)
}

type HierarchyMappingProvider[I any, F FieldProvider] struct {
	name     string
	mapper   chain.Mapper[I, F]
	extender FieldExtender[I, F]
	headers  []string
}

var _ output.MappingProvider[int, FieldProvider] = (*HierarchyMappingProvider[int, FieldProvider])(nil)

func NewHierarchyMappingProvider[I any, F FieldProvider](name string, mapper chain.Mapper[I, F], extender FieldExtender[I, F], headers ...string) *HierarchyMappingProvider[I, F] {
	return &HierarchyMappingProvider[I, F]{name, mapper, extender, slices.Clone(headers)}
}

func (h *HierarchyMappingProvider[I, F]) GetMapping(opts flagutils.OptionSetProvider) (chain.Mapper[I, F], []string, error) {
	c := closure.From[I](opts)
	if c == nil || c.GetExploderFactory(opts) == nil {
		return h.mapper, h.headers, nil
	}
	return func(in I) F {
		return h.extender.Extend(in, h.mapper(in))
	}, append([]string{h.name}, h.headers...), nil
}

////////////////////////////////////////////////////////////////////////////////

func NewStandardHierarchyMappingProvider[I any, F ExtendedFieldProvider](name string, mapper chain.Mapper[I, F], extract chain.Mapper[I, []string], headers ...string) *HierarchyMappingProvider[I, F] {
	return NewHierarchyMappingProvider[I, F](name, mapper, FieldExtenderFunc[I, F](func(o I, in F) F { in.InsertFields(0, extract(o)...); return in }), headers...)
}

func NewTopoHierarchMappingProvider[K comparable, I history.HistoryProvider[K], F ExtendedFieldProvider](name string, mapper chain.Mapper[I, F], headers ...string) *HierarchyMappingProvider[I, F] {
	return NewHierarchyMappingProvider[I, F](name, mapper, FieldExtenderFunc[I, F](func(o I, in F) F { in.InsertFields(0, o.GetHistory().String()); return in }), headers...)
}

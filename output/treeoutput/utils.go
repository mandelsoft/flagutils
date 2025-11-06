package treeoutput

import (
	"github.com/mandelsoft/flagutils/output/tableoutput"
	"github.com/mandelsoft/streaming/chain"
)

func NewHierarchMappingProvider[K, I comparable, O Element[K, I]](name string, mapper chain.Mapper[O, TreeElement[K, I, O]], headers ...string) *tableoutput.HierarchyMappingProvider[O, TreeElement[K, I, O]] {
	f := func(o O, in TreeElement[K, I, O]) TreeElement[K, I, O] {
		in.InsertFields(0, o.GetHistory().String())
		return in
	}
	return tableoutput.NewHierarchyMappingProvider[O, TreeElement[K, I, O]](name, mapper, tableoutput.FieldExtenderFunc[O, TreeElement[K, I, O]](f), headers...)
}

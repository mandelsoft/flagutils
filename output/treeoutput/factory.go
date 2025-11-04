package treeoutput

import (
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/output/tableoutput"
	"github.com/mandelsoft/flagutils/output/treeoutput/topo"
	"github.com/mandelsoft/flagutils/utils/tree"
	"github.com/mandelsoft/goutils/iterutils"
	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/mandelsoft/streaming/chain"
	"slices"
)

// Element described a node element in a tree with
// hierarchy-level names of type K and a unique node identity of
// type I.
type Element[K, I comparable] interface {
	tree.Object[K]
	topo.TopoInfo[K, I]
}

// TreeElement represents an interface combining tree hierarchy, unique
// identity, and field data for processing elements providing access
// to the base object of type O.
type TreeElement[K, I comparable, O Element[K, I]] interface {
	output.FieldProvider
	Element[K, I]
	GetElement() O
}

type element[K, I comparable, O Element[K, I]] struct {
	Element[K, I]
	fields []string
}

func (e *element[K, I, O]) GetElement() O {
	return any(e.Element).(O)
}

func (e *element[K, I, O]) GetFields() []string {
	return e.fields
}

func NewOutputFactory[K, I comparable, O Element[K, I]](opts *TreeOutputOptions[K], cmp topo.ComparerFactory[O], mapper chain.Mapper[O, output.FieldProvider], headers ...string) *tableoutput.OutputFactory[O, TreeElement[K, I, O]] {
	c := chain.Transformed[TreeElement[K, I, O], *tree.TreeObject[K]](treeTransform[K, I, O](cmp))

	return tableoutput.NewExtendedOutputFactory[O, TreeElement[K, I, O]](
		func(o O) TreeElement[K, I, O] {
			return &element[K, I, O]{
				o, mapper(o).GetFields(),
			}
		},
		chain.AddMap[output.FieldProvider](c, treeMapping[K](len(headers), opts)),
		output.ComposeFields(opts.Header(), headers)...,
	)
}

func mapToElement[K, I comparable, O Element[K, I]](e TreeElement[K, I, O]) O {
	return e.GetElement()
}

func treeTransform[K, I comparable, O Element[K, I]](cmp topo.ComparerFactory[O]) func(in []TreeElement[K, I, O]) []*tree.TreeObject[K] {
	return func(in []TreeElement[K, I, O]) []*tree.TreeObject[K] {
		if cmp != nil {
			hcmp := cmp.Comparer(iterutils.ForMapped(mapToElement[K, I, O], in...))
			slices.SortFunc(in, func(a, b TreeElement[K, I, O]) int {
				return hcmp(a.GetElement(), b.GetElement())
			})
		}
		return tree.MapToTree[K](sliceutils.Convert[tree.Object[K]](in), nil)
	}
}

func treeMapping[K comparable](n int, opts *TreeOutputOptions[K]) chain.Mapper[*tree.TreeObject[K], output.FieldProvider] {
	return func(e *tree.TreeObject[K]) output.FieldProvider {
		if e.Object != nil {
			return output.ComposeFields(e.Graph, e.Object)
		}
		return output.ComposeFields(e.Graph+" "+opts.NodeTitle(e), opts.NodeMapping(n, e)) // create empty table line
	}
}

package treeoutput

import (
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/output/tableoutput"
	"github.com/mandelsoft/flagutils/utils/history"
	"github.com/mandelsoft/flagutils/utils/tree"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/mandelsoft/streaming/chain"
	"slices"
)

type Element[K comparable] interface {
	tree.Object[K]
	GetKey() K
}
type TreeElement[K comparable] interface {
	output.FieldProvider
	Element[K]
}

type element[K comparable] struct {
	Element[K]
	fields []string
}

func (e *element[K]) GetFields() []string {
	return e.fields
}

func NewOutputFactory[K comparable, I Element[K]](opts *TreeOutputOptions[K], cmp general.CompareFunc[K], mapper chain.Mapper[I, output.FieldProvider], headers ...string) *tableoutput.OutputFactory[I, TreeElement[K]] {
	c := chain.Transformed[TreeElement[K], *tree.TreeObject[K]](treeTransform[K](cmp))

	return tableoutput.NewExtendedOutputFactory[I, TreeElement[K]](
		func(o I) TreeElement[K] {
			return &element[K]{
				o, mapper(o).GetFields(),
			}
		},
		chain.AddMap[output.FieldProvider](c, treeMapping[K](len(headers), opts)),
		output.ComposeFields(opts.Header(), headers)...,
	)
}

func treeTransform[K comparable](cmp general.CompareFunc[K]) func(in []TreeElement[K]) []*tree.TreeObject[K] {
	return func(in []TreeElement[K]) []*tree.TreeObject[K] {
		hcmp := history.CompareFunc[K](cmp)
		slices.SortFunc(in, func(a, b TreeElement[K]) int {
			return hcmp(a.GetHistory().Add(a.GetKey()), b.GetHistory().Add(b.GetKey()))
		})
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

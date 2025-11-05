package graph

import (
	"fmt"
	"iter"
	"strings"

	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/output/treeoutput/topo"
	"github.com/mandelsoft/flagutils/utils/history"
	"github.com/mandelsoft/flagutils/utils/tree"
	"github.com/mandelsoft/goutils/generics"
)

type Element struct {
	topo.TopoInfo[string, string]
	node *Node
	cont bool
	err  error
}

var _ history.HistoryProvider[string] = (*Element)(nil)
var _ tree.Object[string] = (*Element)(nil)

func NewElement(node *Node, hist history.History[string]) *Element {
	return &Element{node: node, TopoInfo: topo.NewStringIdTopoInfo[string](node.Name(), hist)}
}

func NewErrElement(name string, err error) *Element {
	return &Element{node: nil, TopoInfo: topo.NewStringIdTopoInfo[string](name, nil), err: err}
}

func NewContElement(hist history.History[string], err error) *Element {
	return &Element{node: nil, TopoInfo: topo.NewStringIdTopoInfo[string]("...", hist), cont: true, err: err}
}

func (e *Element) IsNode() *string {
	return generics.PointerTo(e.GetKey())
}

func (e *Element) GetPath() string {
	return strings.Join(e.GetHierarchy(), "/")
}

func (e *Element) GetValue() string {
	if e.node != nil {
		return e.node.Value()
	}
	return ""
}

func (e *Element) AsManifest() any {
	m := map[string]any{}

	m["name"] = e.GetKey()
	m["value"] = e.node.Value()
	if len(e.GetHistory()) > 1 {
		m["path"] = strings.Join(e.GetHistory(), "/")
	}
	if e.err != nil {
		m["error"] = e.err.Error()
	}
	return m
}

////////////////////////////////////////////////////////////////////////////////

type SourceFactory struct {
	graph *Graph
}

func NewSourceFactory(g *Graph) *SourceFactory {
	return &SourceFactory{g}
}

func (s *SourceFactory) Elements(specs output.ElementSpecs) (iter.Seq[*Element], error) {
	list := specs.([]string)
	return func(yield func(*Element) bool) {
		for _, f := range list {
			var e *Element

			n := s.graph.GetRoot(f)
			if n != nil {
				e = NewElement(n, nil)
			} else {
				e = NewErrElement(f, fmt.Errorf("unknown node"))
			}
			if !yield(e) {
				return
			}
		}
	}, nil
}

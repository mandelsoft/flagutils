package graph

import (
	"context"
	"fmt"
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/goutils/set"

	"github.com/mandelsoft/flagutils/utils/history"
	"github.com/mandelsoft/streaming/chain"
)

func ClosureFactory(opts flagutils.OptionSetProvider) chain.ExploderFactory[*Element, *Element] {
	out := output.From[*Element](opts)
	return (&Closure{synthesize: out != nil && out.GetMode() == "tree"}).Create
}

type Closure struct {
	synthesize bool
	found      set.Set[*Node]
}

func (c *Closure) Create(ctx context.Context) chain.Exploder[*Element, *Element] {
	return (&Closure{synthesize: c.synthesize, found: set.Set[*Node]{}}).Closure
}

// Closure creates the transitive closure for a node element
// by recursively following child relations.
func (c *Closure) Closure(e *Element) []*Element {
	if e.err != nil {
		return []*Element{e}
	}

	return c.closure(e.node, nil)
}

func (c *Closure) closure(n *Node, hist history.History[string]) []*Element {
	result := []*Element{NewElement(n, hist)}
	if c.found.Contains(n) && !n.HasChildren() {
		return result
	}
	if !c.synthesize && (hist.Contains(n.Name()) || c.found.Contains(n)) {
		return result
	}
	if hist.Contains(n.Name()) {
		return append(result, NewContElement(hist.Add(n.Name()), fmt.Errorf("cycle")))
	}
	if c.found.Contains(n) {
		return append(result, NewContElement(hist.Add(n.Name()), fmt.Errorf("already shown")))
	}
	hist = hist.Add(n.Name())
	c.found.Add(n)
	for n := range n.Children {
		result = append(result, c.closure(n, hist)...)
	}
	return result
}

package files

import (
	"context"
	"os"
	"sync"

	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/parallel"
	"github.com/mandelsoft/flagutils/utils/pool"
	"github.com/mandelsoft/streaming/chain"
	"github.com/mandelsoft/streaming/processing"
)

// ClosureFactory creates an exploder factory generating
// a sequential or parallel exploder function based on
// flagutils.OptionSetProvider depending
// on the settings of optional parallel.Options.
func ClosureFactory(opts flagutils.OptionSetProvider) chain.ExploderFactory[*Element, *Element] {
	var proc processing.Processing = nil
	p := parallel.From(opts)
	if p != nil {
		proc = p.GetPool()
	}

	// this is the chain.ExploderFactory.
	return func(ctx context.Context) chain.Exploder[*Element, *Element] {
		if proc == nil {
			// use sequential closure calculation chain.Exploder function
			return Closure
		}
		// use parallel closure calculation chain.Exploder function
		return func(in *Element) []*Element {
			result := &result{pool: pool.New[*Element](proc)}
			result.handle(in)
			result.wait()
			return result.result
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

type result struct {
	lock   sync.Mutex
	pool   *pool.Processing[*Element]
	result []*Element
}

func (r *result) wait() {
	r.pool.Wait()
}

func (r *result) Add(e *Element) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.result = append(r.result, e)
}

func (r *result) handle(e *Element) {
	if e.Error != nil || !e.Fi.IsDir() {
		r.Add(e)
		return
	}
	r.pool.Execute(e, r.execute)
}

// execute evaluates a directory.
// For sub-directories a new request is created.
func (r *result) execute(ctx context.Context, p *pool.Processing[*Element], e *Element) {
	entries, err := os.ReadDir(e.GetPath())
	if err != nil {
		e.Error = err
		r.Add(e)
		return
	}
	r.Add(e)
	for _, n := range entries {
		r.handle(NewElement(n.Name(), e.GetHierarchy()))
	}
}

////////////////////////////////////////////////////////////////////////////////

// Closure creates the transitive closure for a file element
// by recursively following directories.
func Closure(e *Element) []*Element {
	result := []*Element{e}
	if e.Error != nil || !e.Fi.IsDir() {
		return result
	}
	entries, err := os.ReadDir(e.GetPath())
	if err != nil {
		e.Error = err
		return result
	}
	for _, n := range entries {
		result = append(result, Closure(NewElement(n.Name(), e.GetHierarchy()))...)
	}
	return result
}

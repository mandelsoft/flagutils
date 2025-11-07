package files

import (
	"context"
	"os"
	"sync"

	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/parallel"
	"github.com/mandelsoft/streaming/chain"
	"github.com/mandelsoft/streaming/processing"
)

func ClosureFactory(opts flagutils.OptionSetProvider) chain.ExploderFactory[*Element, *Element] {
	var pool processing.Processing = nil
	p := parallel.From(opts)
	if p != nil {
		pool = p.GetPool()
	}

	return func(ctx context.Context) chain.Exploder[*Element, *Element] {
		if pool == nil {
			return Closure
		}
		return func(in *Element) []*Element {
			result := &result{pool: pool}
			result.handle(in)
			result.wg.Wait()
			return result.result
		}
	}
}

func (r *result) handle(e *Element) {
	if e.Error != nil || !e.Fi.IsDir() {
		r.Add(e)
		return
	}
	r.wg.Add(1)
	r.pool.Execute(&request{r, e})
}

type result struct {
	lock   sync.Mutex
	pool   processing.Processing
	wg     sync.WaitGroup
	result []*Element
}

func (r *result) Add(e *Element) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.result = append(r.result, e)
}

type request struct {
	r *result
	e *Element
}

var _ processing.Request = (*request)(nil)

// Execute evaluates a directory.
// For sub-directories a new request is created.
func (r *request) Execute(ctx context.Context) {
	defer r.r.wg.Done()
	entries, err := os.ReadDir(r.e.GetPath())
	if err != nil {
		r.e.Error = err
		r.r.Add(r.e)
		return
	}
	r.r.Add(r.e)
	for _, n := range entries {
		r.r.handle(NewElement(n.Name(), r.e.GetHierarchy()))
	}
}

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

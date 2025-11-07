package pool

import (
	"context"
	"fmt"
	"sync"

	"github.com/mandelsoft/streaming/processing"
)

type Processing[E any] struct {
	wg   sync.WaitGroup
	pool processing.Processing
}

func New[E any](pool processing.Processing) *Processing[E] {
	return &Processing[E]{pool: pool}
}

func (p *Processing[E]) Execute(elem E, runner Runner[E]) {
	p.wg.Add(1)
	p.pool.Execute(&request[E]{p, elem, runner})
}

func (p *Processing[E]) Wait() {
	p.wg.Wait()
}

type Runner[E any] func(ctx context.Context, pool *Processing[E], elem E)

type request[E any] struct {
	pool *Processing[E]
	elem E
	f    Runner[E]
}

var _ processing.Request = (*request[int])(nil)

func (r *request[E]) Execute(ctx context.Context) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("%% panic in pool: %v\n", e)
		}
	}()

	r.f(ctx, r.pool, r.elem)
	r.pool.wg.Done()
}

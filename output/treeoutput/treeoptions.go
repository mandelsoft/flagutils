package treeoutput

import "github.com/mandelsoft/flagutils/utils/tree"

type TreeOutputOption[K comparable] interface {
	ApplyTreeOutputOption(*TreeOutputOptions[K])
}

type TreeNodeMappingFunc[K comparable] func(*tree.TreeObject[K]) []string

func (f TreeNodeMappingFunc[K]) ApplyTreeOutputOption(o *TreeOutputOptions[K]) {
	o.nodeMapping = f
}

type TreeNodeTitleFunc[K comparable] func(*tree.TreeObject[K]) string

func (f TreeNodeTitleFunc[K]) ApplyTreeOutputOption(o *TreeOutputOptions[K]) {
	o.nodeTitle = f
}

type TreeOutputOptions[K comparable] struct {
	header      *string
	nodeMapping TreeNodeMappingFunc[K]
	nodeTitle   TreeNodeTitleFunc[K]
}

func (o *TreeOutputOptions[K]) ApplyTreeOutputOption(opts *TreeOutputOptions[K]) {
	if o.nodeMapping != nil {
		opts.nodeMapping = o.nodeMapping
	}
	if o.nodeTitle != nil {
		opts.nodeTitle = o.nodeTitle
	}
	if o.header != nil {
		opts.header = o.header
	}
}

func (o *TreeOutputOptions[K]) WithHeader(name string) *TreeOutputOptions[K] {
	o.header = &name
	return o
}

func (o *TreeOutputOptions[K]) WithNodeTitle(f TreeNodeTitleFunc[K]) *TreeOutputOptions[K] {
	o.nodeTitle = f
	return o
}

func (o *TreeOutputOptions[K]) WithModeMapping(f TreeNodeMappingFunc[K]) *TreeOutputOptions[K] {
	o.nodeMapping = f
	return o
}

func (o *TreeOutputOptions[K]) Apply(opts ...TreeOutputOption[K]) *TreeOutputOptions[K] {
	for _, e := range opts {
		e.ApplyTreeOutputOption(o)
	}
	return o
}

func (o *TreeOutputOptions[K]) NodeMapping(n int, obj *tree.TreeObject[K]) interface{} {
	if o == nil || o.nodeMapping == nil {
		return make([]string, n)
	}
	return o.nodeMapping(obj)
}

func (o *TreeOutputOptions[K]) Header() string {
	if o == nil || o.header == nil {
		return "HIERARCHY"
	}
	return *o.header
}

func (o *TreeOutputOptions[K]) NodeTitle(obj *tree.TreeObject[K]) string {
	if o == nil || o.nodeTitle == nil {
		return obj.Node.String()
	}
	return o.nodeTitle(obj)
}

func WithHeader[K comparable](name string) *TreeOutputOptions[K] {
	return &TreeOutputOptions[K]{header: &name}
}

func WithNodeTitle[K comparable](f TreeNodeTitleFunc[K]) *TreeOutputOptions[K] {
	return &TreeOutputOptions[K]{nodeTitle: f}
}

func WithModeMapping[K comparable](f TreeNodeMappingFunc[K]) *TreeOutputOptions[K] {
	return &TreeOutputOptions[K]{nodeMapping: f}
}

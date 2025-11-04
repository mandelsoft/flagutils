package topo

import "github.com/mandelsoft/goutils/general"

type ComparerFactory[O any] interface {
	Comparer(in []O) general.CompareFunc[O]
}

type ComparerFactoryFunc[O any] func(in []O) general.CompareFunc[O]

func (f ComparerFactoryFunc[O]) Comparer(in []O) general.CompareFunc[O] {
	return f(in)
}

////////////////////////////////////////////////////////////////////////////////

func NewStringIdComparerFactory[T comparable, O TopoInfo[T, string]]() ComparerFactory[O] {
	return NewDefaultComparerFactory[T, string, O](StringIdProviderFunc)
}

func NewStringIdCompareFunc[T comparable, O TopoInfo[T, string]](in []O) general.CompareFunc[O] {
	return NewDefaultCompareFunc(in, StringIdProviderFunc)
}

var _ ComparerFactoryFunc[TopoInfo[string, string]] = NewStringIdCompareFunc[string, TopoInfo[string, string]]

////////////////////////////////////////////////////////////////////////////////

type DefaultComparer[T, I comparable, O TopoInfo[T, I]] struct {
	index  map[I]int
	mapper IdProvider[T, I]
}

// NewDefaultComparerFactory creates a ComparerFactory for a given IdProvider based on
// hierarchy level names of type T and an element identity type I.
// A hierarchy is a set of elements (O) where each element is either a top level element
// or has a parent also member of the element set.
// Such a set can be topologically sorted. The result is an element list, where every element
// is placed later in the list than its parent and earlier than all of its children.
// The order of elements at the same level is not determined. But an element may have
// additional attributes to its hierarchy information. Those attributes may imply
// an order for those elements in the resulting list. Such an order may be given
// by a preordered list of elements according to some desired sorting criteria.
// For example, by ordering the set by other attribute-based compare functions
// (for example, by a sort step of a processing chain.
// This function provides a general.CompareFunc, which can be used by
// a slices.Sort function to provide a topologically sorted element list
// obeying the element order of siblings found in the given initial list.
func NewDefaultComparerFactory[T, I comparable, O TopoInfo[T, I]](mapper IdProvider[T, I]) ComparerFactory[O] {
	return ComparerFactoryFunc[O](func(in []O) general.CompareFunc[O] {
		return NewDefaultCompareFunc(in, mapper)
	})
}

func NewDefaultCompareFunc[T, I comparable, O TopoInfo[T, I]](in []O, mapper IdProvider[T, I]) general.CompareFunc[O] {
	index := make(map[I]int)

	for i, e := range in {
		index[e.GetId()] = i
	}
	return (&DefaultComparer[T, I, O]{
		index:  index,
		mapper: mapper,
	}).Compare
}

func (t *DefaultComparer[T, I, O]) Compare(a, b O) int {
	pa := a.GetHierarchy()
	pb := b.GetHierarchy()
	for i, c := range pa {
		if i >= len(pb) {
			return 1
		}
		if c != pb[i] {
			// compare with level i
			return t.index[t.mapper(pa[:i+1])] - t.index[t.mapper(pb[:i+1])]
		}
	}
	return len(pa) - len(pb)
}

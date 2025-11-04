package topo

import (
	"fmt"
	"github.com/mandelsoft/flagutils/utils/history"
)

type TopoInfo[T, I comparable] interface {
	history.HistoryProvider[T]

	GetHierarchy() []T

	// GetKey provides the name of the actual object.
	GetKey() T

	// GetId provides the globally unique id of the object.
	GetId() I
}

type IdProvider[T, I comparable] func([]T) I

////////////////////////////////////////////////////////////////////////////////

type DefaultTopoInfo[T, I comparable] struct {
	mapper IdProvider[T, I]
	path   history.History[T]
}

func NewDefaultTopoInfo[T, I comparable](mapper IdProvider[T, I], key T, path history.History[T]) *DefaultTopoInfo[T, I] {
	return &DefaultTopoInfo[T, I]{mapper, path.Add(key)}
}

func (t *DefaultTopoInfo[T, I]) GetHierarchy() []T {
	return t.path
}

func (t *DefaultTopoInfo[T, I]) GetHistory() history.History[T] {
	return t.path[:len(t.path)-1]
}

func (t *DefaultTopoInfo[T, I]) GetKey() T {
	return t.path[len(t.path)-1]
}

func (t *DefaultTopoInfo[T, I]) GetId() I {
	return t.mapper(t.path)
}

////////////////////////////////////////////////////////////////////////////////

type StringIdTopoInfo[T comparable] = DefaultTopoInfo[T, string]

func StringIdProviderFunc[T comparable](path []T) string {
	return fmt.Sprint(path)
}

func NewStringIdTopoInfo[T comparable](key T, path history.History[T]) TopoInfo[T, string] {
	return NewDefaultTopoInfo[T, string](StringIdProviderFunc, key, path)
}

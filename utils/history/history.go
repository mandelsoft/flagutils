package history

import (
	"fmt"
	"github.com/mandelsoft/goutils/sliceutils"
)

type HistoryProvider[T comparable] interface {
	GetHistory() History[T]
}

type History[T comparable] []T

func (h History[T]) Add(v ...T) History[T] {
	return sliceutils.CopyAppend(h, v...)
}

func (h History[T]) String() string {
	s := ""
	sep := ""
	for _, e := range h {
		s = fmt.Sprintf("%s%s%s", s, sep, e)
		sep = "->"
	}
	return s
}

func (h History[T]) Contains(c T) bool {
	for _, e := range h {
		if e == c {
			return true
		}
	}
	return false
}

func (h History[T]) HasPrefix(o History[T]) bool {
	if len(o) > len(h) {
		return false
	}
	for i, e := range o {
		if e != h[i] {
			return false
		}
	}
	return true
}

func (h History[T]) Equals(o History[T]) bool {
	if len(h) != len(o) {
		return false
	}
	if h == nil || o == nil {
		return false
	}

	for i, e := range h {
		if e != o[i] {
			return false
		}
	}
	return true
}

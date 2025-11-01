package output

import "github.com/mandelsoft/goutils/sliceutils"

type History[T any] []T

func (h History[T]) Add(v ...T) History[T] {
	return sliceutils.CopyAppend(h, v...)
}

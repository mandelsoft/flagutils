package files

import (
	"os"
)

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
		result = append(result, Closure(NewElement(n.Name(), e.Path))...)
	}
	return result
}

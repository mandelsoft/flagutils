package files

import (
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/goutils/sliceutils"
	"iter"
	"os"
	"strings"
)

type Element struct {
	Name    string
	History output.History[string]
	Error   error
	Fi      os.FileInfo
}

func NewElement(name string, hist output.History[string]) *Element {
	e := &Element{Name: name, History: hist}
	e.Fi, e.Error = os.Stat(e.Path())
	return e
}

func (e *Element) Path() string {
	p := strings.Join(e.History, string(os.PathSeparator))
	if p == "" {
		return e.Name
	}
	return p + string(os.PathSeparator) + e.Name
}

////////////////////////////////////////////////////////////////////////////////

type SourceFactory struct {
	opts *Options
}

func NewSourceFactory(opts flagutils.OptionSetProvider) *SourceFactory {
	return &SourceFactory{From(opts)}
}

func (s *SourceFactory) Elements(specs output.ElementSpecs) (iter.Seq[*Element], error) {
	list := specs.([]string)
	return func(yield func(*Element) bool) {
		for _, f := range list {
			e := NewElement(f, nil)
			if s.opts.dflag || !e.Fi.IsDir() {
				if !yield(e) {
					return
				}
			} else {
				entries, err := os.ReadDir(e.Path())
				if err != nil {
					e.Error = err
					if !yield(e) {
						return
					}
				} else {
					for _, f := range entries {
						e := NewElement(f.Name(), sliceutils.CopyAppend(e.History, e.Name))
						if !yield(e) {
						}
					}
				}
			}
		}
	}, nil
}

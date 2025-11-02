package files

import (
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/closure"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/utils/history"
	"github.com/mandelsoft/flagutils/utils/tree"
	"iter"
	"os"
	"strings"
)

type Element struct {
	Path    []string
	History history.History[string]
	Error   error
	Fi      os.FileInfo
}

var _ history.HistoryProvider[string] = (*Element)(nil)
var _ tree.Object[string] = (*Element)(nil)

func NewElement(name string, hist history.History[string]) *Element {
	p := hist.Add(name)
	e := &Element{Path: p, History: p[:len(p)-1]}
	e.Fi, e.Error = os.Stat(e.GetPath())
	return e
}

func (e *Element) GetKey() string {
	return e.Path[len(e.Path)-1]
}

func (e *Element) IsNode() *string {
	if e.Error == nil && e.Fi.IsDir() {
		return &e.Path[len(e.Path)-1]
	}
	return nil
}

func (e *Element) GetHistory() history.History[string] {
	return e.History
}

func (e *Element) GetPath() string {
	return strings.Join(e.Path, string(os.PathSeparator))
}

func (e *Element) AsManifest() any {
	m := map[string]any{}

	m["name"] = e.GetKey()
	if e.Error != nil {
		m["error"] = e.Error.Error()
	} else {
		if len(e.History) != 0 {
			m["path"] = strings.Join(e.History, string(os.PathSeparator))
		}
		m["mode"] = e.Fi.Mode().String()
		m["fileinfo"] = e.Fi.Mode()
		m["size"] = e.Fi.Size()
		m["modtime"] = e.Fi.ModTime().UnixNano()
	}
	return m
}

////////////////////////////////////////////////////////////////////////////////

type SourceFactory struct {
	opts *Options
}

func NewSourceFactory(opts flagutils.OptionSetProvider) *SourceFactory {
	mine := From(opts)
	all := closure.From[*Element](opts)
	if all != nil && all.GetExploder() != nil {
		mine.dflag = true
	}
	return &SourceFactory{mine}
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
				entries, err := os.ReadDir(e.GetPath())
				if err != nil {
					e.Error = err
					if !yield(e) {
						return
					}
				} else {
					for _, f := range entries {
						e := NewElement(f.Name(), e.Path)
						if !yield(e) {
						}
					}
				}
			}
		}
	}, nil
}

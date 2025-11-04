package files

import (
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/closure"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/output/treeoutput/topo"
	"github.com/mandelsoft/flagutils/utils/history"
	"github.com/mandelsoft/flagutils/utils/tree"
	"github.com/mandelsoft/goutils/generics"
	"iter"
	"os"
	"strings"
)

type Element struct {
	topo.TopoInfo[string, string]
	Error error
	Fi    os.FileInfo
}

var _ history.HistoryProvider[string] = (*Element)(nil)
var _ tree.Object[string] = (*Element)(nil)

func NewElement(name string, hist history.History[string]) *Element {
	e := &Element{TopoInfo: topo.NewStringIdTopoInfo[string](name, hist)}
	e.Fi, e.Error = os.Stat(e.GetPath())
	return e
}

func (e *Element) IsNode() *string {
	if e.Error == nil && e.Fi.IsDir() {
		return generics.PointerTo(e.GetKey())
	}
	return nil
}

func (e *Element) GetPath() string {
	return strings.Join(e.GetHierarchy(), string(os.PathSeparator))
}

func (e *Element) AsManifest() any {
	m := map[string]any{}

	m["name"] = e.GetKey()
	if e.Error != nil {
		m["error"] = e.Error.Error()
	} else {
		if len(e.GetHistory()) > 1 {
			m["path"] = strings.Join(e.GetHistory(), string(os.PathSeparator))
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
						e := NewElement(f.Name(), e.GetHierarchy())
						if !yield(e) {
						}
					}
				}
			}
		}
	}, nil
}

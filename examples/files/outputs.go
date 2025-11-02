package files

import (
	"fmt"
	"os"
	"strings"

	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/output/tableoutput"
	"github.com/mandelsoft/flagutils/output/treeoutput"
)

var OutputsFactory = output.NewOutputsFactory[*Element]().
	Add("", tableoutput.NewOutputFactory[*Element](map_standard, "NAME", "ERROR")).
	Add("wide", tableoutput.NewOutputFactory[*Element](map_wide, "MODE", "NAME", "-SIZE", "ERROR")).
	Add("tree", treeoutput.NewOutputFactory[string, *Element](treeoutput.WithHeader[string](""), strings.Compare, map_tree, "MODE", "NAME", "-SIZE", "ERROR")).
	AddManifestOutputs()

func map_standard(e *Element) output.FieldProvider {
	errstr := ""
	if e.Error != nil {
		errstr = e.Error.Error()
	}
	return output.Fields{e.GetPath(), errstr}
}

func map_wide(e *Element) output.FieldProvider {
	return map_wide_n(e, func(e *Element) string { return e.GetPath() })
}

func map_tree(e *Element) output.FieldProvider {
	return map_wide_n(e, func(e *Element) string { return e.GetKey() })
}

func map_wide_n(e *Element, n func(e *Element) string) output.FieldProvider {
	errstr := ""
	if e.Error != nil {
		errstr = e.Error.Error()
	}

	size := ""
	mode := ""
	if errstr == "" {
		fi, err := os.Stat(e.GetPath())
		if err != nil {
			errstr = err.Error()
		} else {
			size = fmt.Sprintf("%d", fi.Size())
			mode = fi.Mode().String()
		}
	}
	return output.Fields{mode, n(e), size, errstr}
}

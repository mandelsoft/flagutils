package graph

import (
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/output/tableoutput"
	"github.com/mandelsoft/flagutils/output/treeoutput"
	"github.com/mandelsoft/flagutils/output/treeoutput/topo"
)

var OutputsFactory = output.NewOutputsFactory[*Element]().
	Add("", tableoutput.NewOutputFactory[*Element](map_standard, "NAME", "ERROR")).
	Add("wide", tableoutput.NewOutputFactory[*Element](map_wide, "NAME", "VALUE", "ERROR")).
	Add("tree", treeoutput.NewOutputFactory[string, string, *Element](treeoutput.WithHeader[string](""), topo.NewStringIdComparerFactory[string, *Element](), map_tree, "NAME", "VALUE", "ERROR")).
	AddManifestOutputs()

func map_standard(e *Element) output.FieldProvider {
	errstr := ""
	if e.err != nil {
		errstr = e.err.Error()
	}
	return output.Fields{e.GetPath(), errstr}
}

func map_wide(e *Element) output.FieldProvider {
	errstr := ""
	if e.err != nil {
		errstr = e.err.Error()
	}
	return output.Fields{e.GetPath(), e.GetValue(), errstr}
}

func map_tree(e *Element) output.FieldProvider {
	errstr := ""
	if e.err != nil {
		errstr = e.err.Error()
	}
	return output.Fields{e.GetKey(), e.GetValue(), errstr}
}

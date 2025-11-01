package files

import (
	"fmt"
	"os"

	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/output/tableoutput"
)

var OutputsFactory = output.NewOutputsFactory[*Element]().
	Add("", tableoutput.NewOutputFactory[*Element](map_standard, "Name", "Error")).
	Add("wide", tableoutput.NewOutputFactory[*Element](map_wide, "Mode", "Name", "-Size", "Error")).
	AddManifestOutputs()

func map_standard(e *Element) []string {
	err := ""
	if e.Error != nil {
		err = e.Error.Error()
	}
	return []string{e.Path(), err}
}

func map_wide(e *Element) []string {
	errstr := ""
	if e.Error != nil {
		errstr = e.Error.Error()
	}

	size := ""
	mode := ""
	if errstr == "" {
		fi, err := os.Stat(e.Path())
		if err != nil {
			errstr = err.Error()
		} else {
			size = fmt.Sprintf("%d", fi.Size())
			mode = fi.Mode().String()
		}
	}
	return []string{mode, e.Path(), size, errstr}
}

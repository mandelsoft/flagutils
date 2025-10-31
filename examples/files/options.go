package files

import (
	"github.com/mandelsoft/flagutils"
	"github.com/spf13/pflag"
)

type Options struct {
	dflag bool
}

func From(opts flagutils.OptionSetProvider) *Options {
	return flagutils.GetFrom[*Options](opts)
}

func New() *Options {
	return &Options{}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&o.dflag, "directory", "d", false, "show directory instead of files")
}

package tableoutput

import (
	"github.com/mandelsoft/flagutils"
	"github.com/spf13/pflag"
)

func From(opts flagutils.OptionSetProvider) *Options {
	return flagutils.GetFrom[*Options](opts)
}

type Options struct {
	flagutils.OptionBase[*Options]
	optimizedColumns int
	columns          []string
	allColumns       bool
}

func New() *Options {
	o := &Options{}
	o.OptionBase = flagutils.NewBase(o)
	return o
}

func (o *Options) OptimizeColumns(n int) *Options {
	o.optimizedColumns = n
	return o
}

func (o *Options) Columns(c ...string) *Options {
	o.columns = append(o.columns, c...)
	return o
}

func (o *Options) UseAllColumns() bool {
	return o.allColumns
}

func (o *Options) UseColumnOptimization() bool {
	return o.optimizedColumns > 0 && !o.allColumns
}

func (o *Options) GetOptimizedColumns() int {
	if o.UseColumnOptimization() {
		return o.optimizedColumns
	}
	return 0
}

func (o *Options) UseColumns() []string {
	return o.columns
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	if o.optimizedColumns > 0 {
		fs.BoolVarP(&o.allColumns, "all-columns", "", false, "show all table columns")
	}
	fs.StringSliceVarP(&o.columns, o.Long("columns"), o.Short(""), nil, "show selected columns, only")
}

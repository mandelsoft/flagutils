package tableoutput

import (
	"github.com/mandelsoft/flagutils"
	"github.com/spf13/pflag"
)

func From(opts flagutils.OptionSetProvider) *Options {
	return flagutils.GetFrom[*Options](opts)
}

type Options struct {
	optimizedColumns int
	columns          flagutils.SimpleOption[[]string, *Options]
	allColumns       flagutils.SimpleOption[bool, *Options]
}

func New() *Options {
	o := &Options{}
	o.columns = flagutils.NewSimpleOption[[]string](o, nil, "columns", "", "show selected columns")
	o.allColumns = flagutils.NewSimpleOption[bool](o, false, "all-columns", "", "show all table columns")
	return o
}

func (o *Options) WithOptimizedColumns(n int) *Options {
	o.optimizedColumns = n
	return o
}

func (o *Options) WithColumnsNames(long, short string) *Options {
	return o.columns.WithNames(long, short)
}
func (o *Options) WithColumnsDescription(s string) *Options {
	return o.columns.WithDescription(s)
}
func (o *Options) WithAllColumnsNames(long, short string) *Options {
	return o.allColumns.WithNames(long, short)
}
func (o *Options) WithALlColumnsDescription(s string) *Options {
	return o.allColumns.WithDescription(s)
}

func (o *Options) AddColumns(c ...string) *Options {
	o.columns.Set(append(o.columns.Value(), c...))
	return o
}

func (o *Options) UseAllColumns() bool {
	return o.allColumns.Value()
}

func (o *Options) UseColumnOptimization() bool {
	return o.optimizedColumns > 0 && !o.allColumns.Value()
}

func (o *Options) GetOptimizedColumns() int {
	if o.UseColumnOptimization() {
		return o.optimizedColumns
	}
	return 0
}

func (o *Options) UseColumns() []string {
	return o.columns.Value()
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	if o.optimizedColumns > 0 {
		o.allColumns.AddFlags(fs)
	}
	o.columns.AddFlags(fs)
}

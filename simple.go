package flagutils

import "github.com/spf13/pflag"

type VarPFunc[V any] = func(fs *pflag.FlagSet, p *V, name, shorthand string, value V, usage string)

////////////////////////////////////////////////////////////////////////////////

type SimpleOption[V any, T Options] struct {
	self   T
	value  V
	setter VarPFunc[V]
	long   string
	short  string
	desc   string
}

func NewSimpleOption[V any, T Options](self T, def V, long, short, desc string) SimpleOption[V, T] {
	return SimpleOption[V, T]{self: self, setter: VarPFuncFor[V](), value: def, long: long, short: short, desc: desc}
}

func NewSimpleOptionWithSetter[V any, T Options](self T, setter VarPFunc[V], def V, long, short, desc string) SimpleOption[V, T] {
	return SimpleOption[V, T]{self: self, setter: setter, value: def, long: long, short: short, desc: desc}
}

func (o *SimpleOption[V, T]) AddFlags(fs *pflag.FlagSet) {
	o.setter(fs, &o.value, o.long, o.short, o.value, o.desc)
}

func (o *SimpleOption[V, T]) Value() V {
	return o.value
}

func (o *SimpleOption[V, T]) Set(v V) T {
	o.value = v
	return o.self
}

func (o *SimpleOption[V, T]) WithNames(l, s string) T {
	o.long = l
	o.short = s
	return o.self
}

func (o *SimpleOption[V, T]) WithDescription(s string) T {
	o.desc = s
	return o.self
}

////////////////////////////////////////////////////////////////////////////////

func VarPFuncFor[T any]() VarPFunc[T] {
	var v T
	var r any
	switch any(v).(type) {
	case string:
		r = (*pflag.FlagSet).StringVarP
	case int:
		r = (*pflag.FlagSet).IntVarP
	case bool:
		r = (*pflag.FlagSet).BoolVarP
	case []string:
		r = (*pflag.FlagSet).StringSliceVarP
	default:
		panic("unsupported type")
	}
	return r.(VarPFunc[T])
}

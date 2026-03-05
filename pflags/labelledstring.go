package pflags

import (
	"github.com/spf13/pflag"
)

type LabelledString struct {
	Name  string
	Value string
}

type labelledStringValue LabelledString

func NewLabelledStringValue(val LabelledString, p *LabelledString) *labelledStringValue {
	*p = val
	return (*labelledStringValue)(p)
}

func (i *labelledStringValue) String() string {
	if i.Name == "" {
		return ""
	}
	return i.Name + "=" + i.Value
}

func (i *labelledStringValue) Set(s string) error {
	var err error
	i.Name, i.Value, err = parseAssignment(s)
	return err
}

func (i *labelledStringValue) Type() string {
	return "LabelledString"
}

// LabelledStringVar accepts a single key/value pair separated by a =.
// The value is an uninterpreted string (see also LabelledValueVar).
func LabelledStringVar(f *pflag.FlagSet, p *LabelledString, name string, value LabelledString, usage string) {
	f.VarP(NewLabelledStringValue(value, p), name, "", usage)
}

func LabelledStringVarP(flags *pflag.FlagSet, p *LabelledString, name, shorthand string, value LabelledString, usage string) {
	flags.VarP(NewLabelledStringValue(value, p), name, shorthand, usage)
}

func LabelledStringVarPF(flags *pflag.FlagSet, p *LabelledString, name, shorthand string, value LabelledString, usage string) *pflag.Flag {
	return flags.VarPF(NewLabelledStringValue(value, p), name, shorthand, usage)
}

func LabelledStringV(f *pflag.FlagSet, name string, value LabelledString, usage string) *LabelledString {
	p := &LabelledString{}
	LabelledStringVarP(f, p, name, "", value, usage)
	return p
}

func LabelledStringP(f *pflag.FlagSet, name, shorthand string, value LabelledString, usage string) *LabelledString {
	p := &LabelledString{}
	LabelledStringVarP(f, p, name, shorthand, value, usage)
	return p
}

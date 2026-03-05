package pflags

import (
	"encoding/json"

	"github.com/spf13/pflag"
)

type LabelledValue struct {
	Name  string
	Value interface{}
}

type LabelledValueValue LabelledValue

func NewLabelledValueValue(val LabelledValue, p *LabelledValue) *LabelledValueValue {
	*p = val
	return (*LabelledValueValue)(p)
}

func (i *LabelledValueValue) String() string {
	if i.Name == "" {
		return ""
	}
	data, err := json.Marshal(i.Value)
	if err != nil {
		return "error " + err.Error()
	}
	return i.Name + "=" + string(data)
}

func (i *LabelledValueValue) Set(s string) error {
	var err error
	var value string

	i.Name, value, err = parseAssignment(s)
	if err != nil {
		return err
	}
	i.Value, err = parseValue(value)
	return err
}

func (i *LabelledValueValue) Type() string {
	return "<name>=<YAML>"
}

// LabelledValueVar accepts a single key/value pair separated by a =.
// The value is a string evaluated as yaml/json document.
func LabelledValueVar(f *pflag.FlagSet, p *LabelledValue, name string, value LabelledValue, usage string) {
	f.VarP(NewLabelledValueValue(value, p), name, "", usage)
}

func LabelledValueVarP(flags *pflag.FlagSet, p *LabelledValue, name, shorthand string, value LabelledValue, usage string) {
	flags.VarP(NewLabelledValueValue(value, p), name, shorthand, usage)
}

func LabelledValueVarPF(flags *pflag.FlagSet, p *LabelledValue, name, shorthand string, value LabelledValue, usage string) *pflag.Flag {
	return flags.VarPF(NewLabelledValueValue(value, p), name, shorthand, usage)
}

func LabelledValueV(f *pflag.FlagSet, name string, value LabelledValue, usage string) *LabelledValue {
	p := &LabelledValue{}
	LabelledValueVarP(f, p, name, "", value, usage)
	return p
}

func LabelledValueP(f *pflag.FlagSet, name, shorthand string, value LabelledValue, usage string) *LabelledValue {
	p := &LabelledValue{}
	LabelledValueVarP(f, p, name, shorthand, value, usage)
	return p
}

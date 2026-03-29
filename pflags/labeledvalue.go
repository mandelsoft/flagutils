package pflags

import (
	"encoding/json"

	"github.com/spf13/pflag"
)

type LabeledValue struct {
	Name  string
	Value interface{}
}

type LabeledValueValue LabeledValue

func NewLabeledValueValue(val LabeledValue, p *LabeledValue) *LabeledValueValue {
	*p = val
	return (*LabeledValueValue)(p)
}

func (i *LabeledValueValue) String() string {
	if i.Name == "" {
		return ""
	}
	data, err := json.Marshal(i.Value)
	if err != nil {
		return "error " + err.Error()
	}
	return i.Name + "=" + string(data)
}

func (i *LabeledValueValue) Set(s string) error {
	var err error
	var value string

	i.Name, value, err = parseAssignment(s)
	if err != nil {
		return err
	}
	i.Value, err = parseValue(value)
	return err
}

func (i *LabeledValueValue) Type() string {
	return "<name>=<YAML>"
}

func labeledValueConv(sval *LabeledValueValue) (LabeledValue, error) {
	return LabeledValue(*sval), nil
}

// GetLabeledValue return the LabeledValue of a flag with the given name
func GetLabeledValue(f *pflag.FlagSet, name string) (LabeledValue, error) {
	var _nil LabeledValue
	val, err := getFlagType(f, name, "<name>=<YAML>", labeledValueConv)
	if err != nil {
		return _nil, err
	}
	return val, nil
}

// LabeledValueVar accepts a single key/value pair separated by a =.
// The value is a string evaluated as yaml/json document.
func LabeledValueVar(f *pflag.FlagSet, p *LabeledValue, name string, value LabeledValue, usage string) {
	f.VarP(NewLabeledValueValue(value, p), name, "", usage)
}

func LabeledValueVarP(flags *pflag.FlagSet, p *LabeledValue, name, shorthand string, value LabeledValue, usage string) {
	flags.VarP(NewLabeledValueValue(value, p), name, shorthand, usage)
}

func LabeledValueVarPF(flags *pflag.FlagSet, p *LabeledValue, name, shorthand string, value LabeledValue, usage string) *pflag.Flag {
	return flags.VarPF(NewLabeledValueValue(value, p), name, shorthand, usage)
}

func LabeledValueV(f *pflag.FlagSet, name string, value LabeledValue, usage string) *LabeledValue {
	p := &LabeledValue{}
	LabeledValueVarP(f, p, name, "", value, usage)
	return p
}

func LabeledValueP(f *pflag.FlagSet, name, shorthand string, value LabeledValue, usage string) *LabeledValue {
	p := &LabeledValue{}
	LabeledValueVarP(f, p, name, shorthand, value, usage)
	return p
}

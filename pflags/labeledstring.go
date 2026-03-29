package pflags

import (
	"github.com/spf13/pflag"
)

type LabeledString struct {
	Name  string
	Value string
}

type labeledStringValue LabeledString

func NewLabeledStringValue(val LabeledString, p *LabeledString) *labeledStringValue {
	*p = val
	return (*labeledStringValue)(p)
}

func (i *labeledStringValue) String() string {
	if i.Name == "" {
		return ""
	}
	return i.Name + "=" + i.Value
}

func (i *labeledStringValue) Set(s string) error {
	var err error
	i.Name, i.Value, err = parseAssignment(s)
	return err
}

func (i *labeledStringValue) Type() string {
	return "LabeledString"
}

func labeledStringValueConv(sval *labeledStringValue) (LabeledString, error) {
	return LabeledString(*sval), nil
}

// GetLabeledStringValue return the bytes value of a flag with the given name
func GetLabeledStringValue(f *pflag.FlagSet, name string) (LabeledString, error) {
	var _nil LabeledString
	val, err := getFlagType(f, name, "LabeledString", labeledStringValueConv)
	if err != nil {
		return _nil, err
	}
	return val, nil
}

// LabeledStringVar accepts a single key/value pair separated by a =.
// The value is an uninterpreted string (see also LabeledValueVar).
func LabeledStringVar(f *pflag.FlagSet, p *LabeledString, name string, value LabeledString, usage string) {
	f.VarP(NewLabeledStringValue(value, p), name, "", usage)
}

func LabeledStringVarP(flags *pflag.FlagSet, p *LabeledString, name, shorthand string, value LabeledString, usage string) {
	flags.VarP(NewLabeledStringValue(value, p), name, shorthand, usage)
}

func LabeledStringVarPF(flags *pflag.FlagSet, p *LabeledString, name, shorthand string, value LabeledString, usage string) *pflag.Flag {
	return flags.VarPF(NewLabeledStringValue(value, p), name, shorthand, usage)
}

func LabeledStringV(f *pflag.FlagSet, name string, value LabeledString, usage string) *LabeledString {
	p := &LabeledString{}
	LabeledStringVarP(f, p, name, "", value, usage)
	return p
}

func LabeledStringP(f *pflag.FlagSet, name, shorthand string, value LabeledString, usage string) *LabeledString {
	p := &LabeledString{}
	LabeledStringVarP(f, p, name, shorthand, value, usage)
	return p
}

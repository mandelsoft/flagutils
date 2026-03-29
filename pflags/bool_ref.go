package pflags

import (
	"strconv"

	"github.com/spf13/pflag"
)

// -- bool Ref Value
type boolRefValue refValue[bool]

func newBoolRefValue(val *bool, p **bool) *boolRefValue {
	*p = val
	return &boolRefValue{ptr: p}
}

func (s *boolRefValue) Set(val string) error {
	v, err := strconv.ParseBool(val)
	*s.ptr = &v
	return err
}

func (s *boolRefValue) Type() string {
	return "*bool"
}

func (b *boolRefValue) IsBoolFlag() bool { return true }

func (s *boolRefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatBool(**s.ptr)
}

func boolRefConv(sval *boolRefValue) (*bool, error) {
	return *sval.ptr, nil
}

// GetBoolRef return the bool ref value of a flag with the given name
func GetBoolRef(f *pflag.FlagSet, name string) (*bool, error) {

	val, err := getFlagType(f, name, "*bool", boolRefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// BoolRefVar defines a bool ref flag with specified name, default value, and usage bool.
// The argument p points to a bool pointer variable in which to store the value of the flag as referenced value.
func BoolRefVar(f *pflag.FlagSet, p **bool, name string, value *bool, usage string) {
	f.VarP(newBoolRefValue(value, p), name, "", usage)
}

// BoolRefVarP is like BoolRefVar, but accepts a shorthand letter that can be used after a single dash.
func BoolRefVarP(f *pflag.FlagSet, p **bool, name, shorthand string, value *bool, usage string) {
	f.VarP(newBoolRefValue(value, p), name, shorthand, usage)
}

// BoolRef defines a bool ref flag with specified name, default value, and usage bool.
// The return value is the address of a bool pointer variable that stores the value of the flag.
func BoolRef(f *pflag.FlagSet, name string, value *bool, usage string) **bool {
	p := new(*bool)
	BoolRefVarP(f, p, name, "", value, usage)
	return p
}

// BoolRefP is like BoolRef, but accepts a shorthand letter that can be used after a single dash.
func BoolRefP(f *pflag.FlagSet, name, shorthand string, value *bool, usage string) **bool {
	p := new(*bool)
	BoolRefVarP(f, p, name, shorthand, value, usage)
	return p
}

// BoolRefVarPF is like BoolRefVarP, but returns the created flag.
func BoolRefVarPF(f *pflag.FlagSet, p **bool, name, shorthand string, value *bool, usage string) *pflag.Flag {
	return f.VarPF(newBoolRefValue(value, p), name, shorthand, usage)
}

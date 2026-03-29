package pflags

import (
	"strconv"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- int16 Ref Value
type int16RefValue refValue[int16]

func newInt16RefValue(val *int16, p **int16) *int16RefValue {
	*p = val
	return &int16RefValue{ptr: p}
}

func (s *int16RefValue) Set(val string) error {
	v, err := strconv.ParseInt(val, 10, 16)
	*s.ptr = generics.PointerTo(int16(v))
	return err
}

func (s *int16RefValue) Type() string {
	return "*int16"
}

func (s *int16RefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatInt(int64(**s.ptr), 10)
}

func int16RefConv(sval *int16RefValue) (*int16, error) {
	return *sval.ptr, nil
}

// GetInt16Ref return the int16 ref value of a flag with the given name
func GetInt16Ref(f *pflag.FlagSet, name string) (*int16, error) {

	val, err := getFlagType(f, name, "*int16", int16RefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Int16RefVar defines a int16 ref flag with specified name, default value, and usage int16.
// The argument p points to a int16 pointer variable in which to store the value of the flag as referenced value.
func Int16RefVar(f *pflag.FlagSet, p **int16, name string, value *int16, usage string) {
	f.VarP(newInt16RefValue(value, p), name, "", usage)
}

// Int16RefVarP is like Int16RefVar, but accepts a shorthand letter that can be used after a single dash.
func Int16RefVarP(f *pflag.FlagSet, p **int16, name, shorthand string, value *int16, usage string) {
	f.VarP(newInt16RefValue(value, p), name, shorthand, usage)
}

// Int16Ref defines a int16 ref flag with specified name, default value, and usage int16.
// The return value is the address of a int16 pointer variable that stores the value of the flag.
func Int16Ref(f *pflag.FlagSet, name string, value *int16, usage string) **int16 {
	p := new(*int16)
	Int16RefVarP(f, p, name, "", value, usage)
	return p
}

// Int16RefP is like Int16Ref, but accepts a shorthand letter that can be used after a single dash.
func Int16RefP(f *pflag.FlagSet, name, shorthand string, value *int16, usage string) **int16 {
	p := new(*int16)
	Int16RefVarP(f, p, name, shorthand, value, usage)
	return p
}

// Int16RefVarPF is like Int16RefVarP, but returns the created flag.
func Int16RefVarPF(f *pflag.FlagSet, p **int16, name, shorthand string, value *int16, usage string) *pflag.Flag {
	return f.VarPF(newInt16RefValue(value, p), name, shorthand, usage)
}

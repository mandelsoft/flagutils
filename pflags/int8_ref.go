package pflags

import (
	"strconv"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- int8 Ref Value
type int8RefValue refValue[int8]

func newInt8RefValue(val *int8, p **int8) *int8RefValue {
	*p = val
	return &int8RefValue{ptr: p}
}

func (s *int8RefValue) Set(val string) error {
	v, err := strconv.ParseInt(val, 10, 8)
	*s.ptr = generics.PointerTo(int8(v))
	return err
}

func (s *int8RefValue) Type() string {
	return "*int8"
}

func (s *int8RefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatInt(int64(**s.ptr), 10)
}

func int8RefConv(sval *int8RefValue) (*int8, error) {
	return *sval.ptr, nil
}

// GetInt8Ref return the int8 ref value of a flag with the given name
func GetInt8Ref(f *pflag.FlagSet, name string) (*int8, error) {

	val, err := getFlagType(f, name, "*int8", int8RefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Int8RefVar defines a int8 ref flag with specified name, default value, and usage int8.
// The argument p points to a int8 pointer variable in which to store the value of the flag as referenced value.
func Int8RefVar(f *pflag.FlagSet, p **int8, name string, value *int8, usage string) {
	f.VarP(newInt8RefValue(value, p), name, "", usage)
}

// Int8RefVarP is like Int8RefVar, but accepts a shorthand letter that can be used after a single dash.
func Int8RefVarP(f *pflag.FlagSet, p **int8, name, shorthand string, value *int8, usage string) {
	f.VarP(newInt8RefValue(value, p), name, shorthand, usage)
}

// Int8Ref defines a int8 ref flag with specified name, default value, and usage int8.
// The return value is the address of a int8 pointer variable that stores the value of the flag.
func Int8Ref(f *pflag.FlagSet, name string, value *int8, usage string) **int8 {
	p := new(*int8)
	Int8RefVarP(f, p, name, "", value, usage)
	return p
}

// Int8RefP is like Int8Ref, but accepts a shorthand letter that can be used after a single dash.
func Int8RefP(f *pflag.FlagSet, name, shorthand string, value *int8, usage string) **int8 {
	p := new(*int8)
	Int8RefVarP(f, p, name, shorthand, value, usage)
	return p
}

// Int8RefVarPF is like Int8RefVarP, but returns the created flag.
func Int8RefVarPF(f *pflag.FlagSet, p **int8, name, shorthand string, value *int8, usage string) *pflag.Flag {
	return f.VarPF(newInt8RefValue(value, p), name, shorthand, usage)
}

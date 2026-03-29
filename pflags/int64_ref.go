package pflags

import (
	"strconv"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- int64 Ref Value
type int64RefValue refValue[int64]

func newInt64RefValue(val *int64, p **int64) *int64RefValue {
	*p = val
	return &int64RefValue{ptr: p}
}

func (s *int64RefValue) Set(val string) error {
	v, err := strconv.ParseInt(val, 10, 64)
	*s.ptr = generics.PointerTo(int64(v))
	return err
}

func (s *int64RefValue) Type() string {
	return "*int64"
}

func (s *int64RefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatInt(int64(**s.ptr), 10)
}

func int64RefConv(sval *int64RefValue) (*int64, error) {
	return *sval.ptr, nil
}

// GetInt64Ref return the int64 ref value of a flag with the given name
func GetInt64Ref(f *pflag.FlagSet, name string) (*int64, error) {

	val, err := getFlagType(f, name, "*int64", int64RefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Int64RefVar defines a int64 ref flag with specified name, default value, and usage int64.
// The argument p points to a int64 pointer variable in which to store the value of the flag as referenced value.
func Int64RefVar(f *pflag.FlagSet, p **int64, name string, value *int64, usage string) {
	f.VarP(newInt64RefValue(value, p), name, "", usage)
}

// Int64RefVarP is like Int64RefVar, but accepts a shorthand letter that can be used after a single dash.
func Int64RefVarP(f *pflag.FlagSet, p **int64, name, shorthand string, value *int64, usage string) {
	f.VarP(newInt64RefValue(value, p), name, shorthand, usage)
}

// Int64Ref defines a int64 ref flag with specified name, default value, and usage int64.
// The return value is the address of a int64 pointer variable that stores the value of the flag.
func Int64Ref(f *pflag.FlagSet, name string, value *int64, usage string) **int64 {
	p := new(*int64)
	Int64RefVarP(f, p, name, "", value, usage)
	return p
}

// Int64RefP is like Int64Ref, but accepts a shorthand letter that can be used after a single dash.
func Int64RefP(f *pflag.FlagSet, name, shorthand string, value *int64, usage string) **int64 {
	p := new(*int64)
	Int64RefVarP(f, p, name, shorthand, value, usage)
	return p
}

// Int64RefVarPF is like Int64RefVarP, but returns the created flag.
func Int64RefVarPF(f *pflag.FlagSet, p **int64, name, shorthand string, value *int64, usage string) *pflag.Flag {
	return f.VarPF(newInt64RefValue(value, p), name, shorthand, usage)
}

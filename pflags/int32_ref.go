package pflags

import (
	"strconv"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- int32 Ref Value
type int32RefValue refValue[int32]

func newInt32RefValue(val *int32, p **int32) *int32RefValue {
	*p = val
	return &int32RefValue{ptr: p}
}

func (s *int32RefValue) Set(val string) error {
	v, err := strconv.ParseInt(val, 10, 32)
	*s.ptr = generics.PointerTo(int32(v))
	return err
}

func (s *int32RefValue) Type() string {
	return "*int32"
}

func (s *int32RefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatInt(int64(**s.ptr), 10)
}

func int32RefConv(sval *int32RefValue) (*int32, error) {
	return *sval.ptr, nil
}

// GetInt32Ref return the int32 ref value of a flag with the given name
func GetInt32Ref(f *pflag.FlagSet, name string) (*int32, error) {

	val, err := getFlagType(f, name, "*int32", int32RefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Int32RefVar defines a int32 ref flag with specified name, default value, and usage int32.
// The argument p points to a int32 pointer variable in which to store the value of the flag as referenced value.
func Int32RefVar(f *pflag.FlagSet, p **int32, name string, value *int32, usage string) {
	f.VarP(newInt32RefValue(value, p), name, "", usage)
}

// Int32RefVarP is like Int32RefVar, but accepts a shorthand letter that can be used after a single dash.
func Int32RefVarP(f *pflag.FlagSet, p **int32, name, shorthand string, value *int32, usage string) {
	f.VarP(newInt32RefValue(value, p), name, shorthand, usage)
}

// Int32Ref defines a int32 ref flag with specified name, default value, and usage int32.
// The return value is the address of a int32 pointer variable that stores the value of the flag.
func Int32Ref(f *pflag.FlagSet, name string, value *int32, usage string) **int32 {
	p := new(*int32)
	Int32RefVarP(f, p, name, "", value, usage)
	return p
}

// Int32RefP is like Int32Ref, but accepts a shorthand letter that can be used after a single dash.
func Int32RefP(f *pflag.FlagSet, name, shorthand string, value *int32, usage string) **int32 {
	p := new(*int32)
	Int32RefVarP(f, p, name, shorthand, value, usage)
	return p
}

// Int32RefVarPF is like Int32RefVarP, but returns the created flag.
func Int32RefVarPF(f *pflag.FlagSet, p **int32, name, shorthand string, value *int32, usage string) *pflag.Flag {
	return f.VarPF(newInt32RefValue(value, p), name, shorthand, usage)
}

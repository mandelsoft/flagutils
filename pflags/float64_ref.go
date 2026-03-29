package pflags

import (
	"strconv"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- float64 Ref Value
type float64RefValue refValue[float64]

func newFloat64RefValue(val *float64, p **float64) *float64RefValue {
	*p = val
	return &float64RefValue{ptr: p}
}

func (s *float64RefValue) Set(val string) error {
	v, err := strconv.ParseFloat(val, 64)
	*s.ptr = generics.PointerTo(float64(v))
	return err
}

func (s *float64RefValue) Type() string {
	return "*float64"
}

func (s *float64RefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatFloat(float64(**s.ptr), 'g', -1, 64)
}

func float64RefConv(sval *float64RefValue) (*float64, error) {
	return *sval.ptr, nil
}

// GetFloat64Ref return the float64 ref value of a flag with the given name
func GetFloat64Ref(f *pflag.FlagSet, name string) (*float64, error) {

	val, err := getFlagType(f, name, "*float64", float64RefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Float64RefVar defines a float64 ref flag with specified name, default value, and usage float64.
// The argument p points to a float64 pointer variable in which to store the value of the flag as referenced value.
func Float64RefVar(f *pflag.FlagSet, p **float64, name string, value *float64, usage string) {
	f.VarP(newFloat64RefValue(value, p), name, "", usage)
}

// Float64RefVarP is like Float64RefVar, but accepts a shorthand letter that can be used after a single dash.
func Float64RefVarP(f *pflag.FlagSet, p **float64, name, shorthand string, value *float64, usage string) {
	f.VarP(newFloat64RefValue(value, p), name, shorthand, usage)
}

// Float64Ref defines a float64 ref flag with specified name, default value, and usage float64.
// The return value is the address of a float64 pointer variable that stores the value of the flag.
func Float64Ref(f *pflag.FlagSet, name string, value *float64, usage string) **float64 {
	p := new(*float64)
	Float64RefVarP(f, p, name, "", value, usage)
	return p
}

// Float64RefP is like Float64Ref, but accepts a shorthand letter that can be used after a single dash.
func Float64RefP(f *pflag.FlagSet, name, shorthand string, value *float64, usage string) **float64 {
	p := new(*float64)
	Float64RefVarP(f, p, name, shorthand, value, usage)
	return p
}

// Float64RefVarPF is like Float64RefVarP, but returns the created flag.
func Float64RefVarPF(f *pflag.FlagSet, p **float64, name, shorthand string, value *float64, usage string) *pflag.Flag {
	return f.VarPF(newFloat64RefValue(value, p), name, shorthand, usage)
}

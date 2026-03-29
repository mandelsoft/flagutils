package pflags

import (
	"strconv"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- float32 Ref Value
type float32RefValue refValue[float32]

func newFloat32RefValue(val *float32, p **float32) *float32RefValue {
	*p = val
	return &float32RefValue{ptr: p}
}

func (s *float32RefValue) Set(val string) error {
	v, err := strconv.ParseFloat(val, 32)
	*s.ptr = generics.PointerTo(float32(v))
	return err
}

func (s *float32RefValue) Type() string {
	return "*float32"
}

func (s *float32RefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatFloat(float64(**s.ptr), 'g', -1, 32)
}

func float32RefConv(sval *float32RefValue) (*float32, error) {
	return *sval.ptr, nil
}

// GetFloat32Ref return the float32 ref value of a flag with the given name
func GetFloat32Ref(f *pflag.FlagSet, name string) (*float32, error) {

	val, err := getFlagType(f, name, "*float32", float32RefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Float32RefVar defines a float32 ref flag with specified name, default value, and usage float32.
// The argument p points to a float32 pointer variable in which to store the value of the flag as referenced value.
func Float32RefVar(f *pflag.FlagSet, p **float32, name string, value *float32, usage string) {
	f.VarP(newFloat32RefValue(value, p), name, "", usage)
}

// Float32RefVarP is like Float32RefVar, but accepts a shorthand letter that can be used after a single dash.
func Float32RefVarP(f *pflag.FlagSet, p **float32, name, shorthand string, value *float32, usage string) {
	f.VarP(newFloat32RefValue(value, p), name, shorthand, usage)
}

// Float32Ref defines a float32 ref flag with specified name, default value, and usage float32.
// The return value is the address of a float32 pointer variable that stores the value of the flag.
func Float32Ref(f *pflag.FlagSet, name string, value *float32, usage string) **float32 {
	p := new(*float32)
	Float32RefVarP(f, p, name, "", value, usage)
	return p
}

// Float32RefP is like Float32Ref, but accepts a shorthand letter that can be used after a single dash.
func Float32RefP(f *pflag.FlagSet, name, shorthand string, value *float32, usage string) **float32 {
	p := new(*float32)
	Float32RefVarP(f, p, name, shorthand, value, usage)
	return p
}

// Float32RefVarPF is like Float32RefVarP, but returns the created flag.
func Float32RefVarPF(f *pflag.FlagSet, p **float32, name, shorthand string, value *float32, usage string) *pflag.Flag {
	return f.VarPF(newFloat32RefValue(value, p), name, shorthand, usage)
}

package pflags

import (
	"strconv"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- uint32 Ref Value
type uint32RefValue refValue[uint32]

func newUint32RefValue(val *uint32, p **uint32) *uint32RefValue {
	*p = val
	return &uint32RefValue{ptr: p}
}

func (s *uint32RefValue) Set(val string) error {
	v, err := strconv.ParseUint(val, 10, 32)
	*s.ptr = generics.PointerTo(uint32(v))
	return err
}

func (s *uint32RefValue) Type() string {
	return "*uint32"
}

func (s *uint32RefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatUint(uint64(**s.ptr), 10)
}

func uint32RefConv(sval *uint32RefValue) (*uint32, error) {
	return *sval.ptr, nil
}

// GetUint32Ref return the uint32 ref value of a flag with the given name
func GetUint32Ref(f *pflag.FlagSet, name string) (*uint32, error) {

	val, err := getFlagType(f, name, "*uint32", uint32RefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Uint32RefVar defines a uint32 ref flag with specified name, default value, and usage uint32.
// The argument p points to a uint32 pointer variable in which to store the value of the flag as referenced value.
func Uint32RefVar(f *pflag.FlagSet, p **uint32, name string, value *uint32, usage string) {
	f.VarP(newUint32RefValue(value, p), name, "", usage)
}

// Uint32RefVarP is like Uint32RefVar, but accepts a shorthand letter that can be used after a single dash.
func Uint32RefVarP(f *pflag.FlagSet, p **uint32, name, shorthand string, value *uint32, usage string) {
	f.VarP(newUint32RefValue(value, p), name, shorthand, usage)
}

// Uint32Ref defines a uint32 ref flag with specified name, default value, and usage uint32.
// The return value is the address of a uint32 pointer variable that stores the value of the flag.
func Uint32Ref(f *pflag.FlagSet, name string, value *uint32, usage string) **uint32 {
	p := new(*uint32)
	Uint32RefVarP(f, p, name, "", value, usage)
	return p
}

// Uint32RefP is like Uint32Ref, but accepts a shorthand letter that can be used after a single dash.
func Uint32RefP(f *pflag.FlagSet, name, shorthand string, value *uint32, usage string) **uint32 {
	p := new(*uint32)
	Uint32RefVarP(f, p, name, shorthand, value, usage)
	return p
}

// Uint32RefVarPF is like Uint32RefVarP, but returns the created flag.
func Uint32RefVarPF(f *pflag.FlagSet, p **uint32, name, shorthand string, value *uint32, usage string) *pflag.Flag {
	return f.VarPF(newUint32RefValue(value, p), name, shorthand, usage)
}

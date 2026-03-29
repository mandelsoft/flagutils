package pflags

import (
	"strconv"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- uint64 Ref Value
type uint64RefValue refValue[uint64]

func newUint64RefValue(val *uint64, p **uint64) *uint64RefValue {
	*p = val
	return &uint64RefValue{ptr: p}
}

func (s *uint64RefValue) Set(val string) error {
	v, err := strconv.ParseUint(val, 10, 64)
	*s.ptr = generics.PointerTo(uint64(v))
	return err
}

func (s *uint64RefValue) Type() string {
	return "*uint64"
}

func (s *uint64RefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatUint(uint64(**s.ptr), 10)
}

func uint64RefConv(sval *uint64RefValue) (*uint64, error) {
	return *sval.ptr, nil
}

// GetUint64Ref return the uint64 ref value of a flag with the given name
func GetUint64Ref(f *pflag.FlagSet, name string) (*uint64, error) {

	val, err := getFlagType(f, name, "*uint64", uint64RefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Uint64RefVar defines a uint64 ref flag with specified name, default value, and usage uint64.
// The argument p points to a uint64 pointer variable in which to store the value of the flag as referenced value.
func Uint64RefVar(f *pflag.FlagSet, p **uint64, name string, value *uint64, usage string) {
	f.VarP(newUint64RefValue(value, p), name, "", usage)
}

// Uint64RefVarP is like Uint64RefVar, but accepts a shorthand letter that can be used after a single dash.
func Uint64RefVarP(f *pflag.FlagSet, p **uint64, name, shorthand string, value *uint64, usage string) {
	f.VarP(newUint64RefValue(value, p), name, shorthand, usage)
}

// Uint64Ref defines a uint64 ref flag with specified name, default value, and usage uint64.
// The return value is the address of a uint64 pointer variable that stores the value of the flag.
func Uint64Ref(f *pflag.FlagSet, name string, value *uint64, usage string) **uint64 {
	p := new(*uint64)
	Uint64RefVarP(f, p, name, "", value, usage)
	return p
}

// Uint64RefP is like Uint64Ref, but accepts a shorthand letter that can be used after a single dash.
func Uint64RefP(f *pflag.FlagSet, name, shorthand string, value *uint64, usage string) **uint64 {
	p := new(*uint64)
	Uint64RefVarP(f, p, name, shorthand, value, usage)
	return p
}

// Uint64RefVarPF is like Uint64RefVarP, but returns the created flag.
func Uint64RefVarPF(f *pflag.FlagSet, p **uint64, name, shorthand string, value *uint64, usage string) *pflag.Flag {
	return f.VarPF(newUint64RefValue(value, p), name, shorthand, usage)
}

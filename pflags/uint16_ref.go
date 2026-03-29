package pflags

import (
	"strconv"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- uint16 Ref Value
type uint16RefValue refValue[uint16]

func newUint16RefValue(val *uint16, p **uint16) *uint16RefValue {
	*p = val
	return &uint16RefValue{ptr: p}
}

func (s *uint16RefValue) Set(val string) error {
	v, err := strconv.ParseUint(val, 10, 16)
	*s.ptr = generics.PointerTo(uint16(v))
	return err
}

func (s *uint16RefValue) Type() string {
	return "*uint16"
}

func (s *uint16RefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatUint(uint64(**s.ptr), 10)
}

func uint16RefConv(sval *uint16RefValue) (*uint16, error) {
	return *sval.ptr, nil
}

// GetUint16Ref return the uint16 ref value of a flag with the given name
func GetUint16Ref(f *pflag.FlagSet, name string) (*uint16, error) {

	val, err := getFlagType(f, name, "*uint16", uint16RefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Uint16RefVar defines a uint16 ref flag with specified name, default value, and usage uint16.
// The argument p points to a uint16 pointer variable in which to store the value of the flag as referenced value.
func Uint16RefVar(f *pflag.FlagSet, p **uint16, name string, value *uint16, usage string) {
	f.VarP(newUint16RefValue(value, p), name, "", usage)
}

// Uint16RefVarP is like Uint16RefVar, but accepts a shorthand letter that can be used after a single dash.
func Uint16RefVarP(f *pflag.FlagSet, p **uint16, name, shorthand string, value *uint16, usage string) {
	f.VarP(newUint16RefValue(value, p), name, shorthand, usage)
}

// Uint16Ref defines a uint16 ref flag with specified name, default value, and usage uint16.
// The return value is the address of a uint16 pointer variable that stores the value of the flag.
func Uint16Ref(f *pflag.FlagSet, name string, value *uint16, usage string) **uint16 {
	p := new(*uint16)
	Uint16RefVarP(f, p, name, "", value, usage)
	return p
}

// Uint16RefP is like Uint16Ref, but accepts a shorthand letter that can be used after a single dash.
func Uint16RefP(f *pflag.FlagSet, name, shorthand string, value *uint16, usage string) **uint16 {
	p := new(*uint16)
	Uint16RefVarP(f, p, name, shorthand, value, usage)
	return p
}

// Uint16RefVarPF is like Uint16RefVarP, but returns the created flag.
func Uint16RefVarPF(f *pflag.FlagSet, p **uint16, name, shorthand string, value *uint16, usage string) *pflag.Flag {
	return f.VarPF(newUint16RefValue(value, p), name, shorthand, usage)
}

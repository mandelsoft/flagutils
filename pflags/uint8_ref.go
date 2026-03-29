package pflags

import (
	"strconv"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- uint8 Ref Value
type uint8RefValue refValue[uint8]

func newUint8RefValue(val *uint8, p **uint8) *uint8RefValue {
	*p = val
	return &uint8RefValue{ptr: p}
}

func (s *uint8RefValue) Set(val string) error {
	v, err := strconv.ParseUint(val, 10, 8)
	*s.ptr = generics.PointerTo(uint8(v))
	return err
}

func (s *uint8RefValue) Type() string {
	return "*uint8"
}

func (s *uint8RefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatUint(uint64(**s.ptr), 10)
}

func uint8RefConv(sval *uint8RefValue) (*uint8, error) {
	return *sval.ptr, nil
}

// GetUint8Ref return the uint8 ref value of a flag with the given name
func GetUint8Ref(f *pflag.FlagSet, name string) (*uint8, error) {

	val, err := getFlagType(f, name, "*uint8", uint8RefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Uint8RefVar defines a uint8 ref flag with specified name, default value, and usage uint8.
// The argument p points to a uint8 pointer variable in which to store the value of the flag as referenced value.
func Uint8RefVar(f *pflag.FlagSet, p **uint8, name string, value *uint8, usage string) {
	f.VarP(newUint8RefValue(value, p), name, "", usage)
}

// Uint8RefVarP is like Uint8RefVar, but accepts a shorthand letter that can be used after a single dash.
func Uint8RefVarP(f *pflag.FlagSet, p **uint8, name, shorthand string, value *uint8, usage string) {
	f.VarP(newUint8RefValue(value, p), name, shorthand, usage)
}

// Uint8Ref defines a uint8 ref flag with specified name, default value, and usage uint8.
// The return value is the address of a uint8 pointer variable that stores the value of the flag.
func Uint8Ref(f *pflag.FlagSet, name string, value *uint8, usage string) **uint8 {
	p := new(*uint8)
	Uint8RefVarP(f, p, name, "", value, usage)
	return p
}

// Uint8RefP is like Uint8Ref, but accepts a shorthand letter that can be used after a single dash.
func Uint8RefP(f *pflag.FlagSet, name, shorthand string, value *uint8, usage string) **uint8 {
	p := new(*uint8)
	Uint8RefVarP(f, p, name, shorthand, value, usage)
	return p
}

// Uint8RefVarPF is like Uint8RefVarP, but returns the created flag.
func Uint8RefVarPF(f *pflag.FlagSet, p **uint8, name, shorthand string, value *uint8, usage string) *pflag.Flag {
	return f.VarPF(newUint8RefValue(value, p), name, shorthand, usage)
}

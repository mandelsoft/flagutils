package pflags

import (
	"strconv"
	"unsafe"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- uint Ref Value
type uintRefValue refValue[uint]

func newUintRefValue(val *uint, p **uint) *uintRefValue {
	*p = val
	return &uintRefValue{ptr: p}
}

func (s *uintRefValue) Set(val string) error {
	v, err := strconv.ParseUint(val, 10, int(unsafe.Sizeof(uint(0))))
	*s.ptr = generics.PointerTo(uint(v))
	return err
}

func (s *uintRefValue) Type() string {
	return "*uint"
}

func (s *uintRefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatUint(uint64(**s.ptr), 10)
}

func uintRefConv(sval *uintRefValue) (*uint, error) {
	return *sval.ptr, nil
}

// GetUintRef return the uint ref value of a flag with the given name
func GetUintRef(f *pflag.FlagSet, name string) (*uint, error) {

	val, err := getFlagType(f, name, "*uint", uintRefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// UintRefVar defines a uint ref flag with specified name, default value, and usage uint.
// The argument p points to a uint pointer variable in which to store the value of the flag as referenced value.
func UintRefVar(f *pflag.FlagSet, p **uint, name string, value *uint, usage string) {
	f.VarP(newUintRefValue(value, p), name, "", usage)
}

// UintRefVarP is like UintRefVar, but accepts a shorthand letter that can be used after a single dash.
func UintRefVarP(f *pflag.FlagSet, p **uint, name, shorthand string, value *uint, usage string) {
	f.VarP(newUintRefValue(value, p), name, shorthand, usage)
}

// UintRef defines a uint ref flag with specified name, default value, and usage uint.
// The return value is the address of a uint pointer variable that stores the value of the flag.
func UintRef(f *pflag.FlagSet, name string, value *uint, usage string) **uint {
	p := new(*uint)
	UintRefVarP(f, p, name, "", value, usage)
	return p
}

// UintRefP is like UintRef, but accepts a shorthand letter that can be used after a single dash.
func UintRefP(f *pflag.FlagSet, name, shorthand string, value *uint, usage string) **uint {
	p := new(*uint)
	UintRefVarP(f, p, name, shorthand, value, usage)
	return p
}

// UintRefVarPF is like UintRefVarP, but returns the created flag.
func UintRefVarPF(f *pflag.FlagSet, p **uint, name, shorthand string, value *uint, usage string) *pflag.Flag {
	return f.VarPF(newUintRefValue(value, p), name, shorthand, usage)
}

package pflags

import (
	"strconv"
	"unsafe"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// -- int Ref Value
type intRefValue refValue[int]

func newIntRefValue(val *int, p **int) *intRefValue {
	*p = val
	return &intRefValue{ptr: p}
}

func (s *intRefValue) Set(val string) error {
	v, err := strconv.ParseInt(val, 10, int(unsafe.Sizeof(int(0)))*8)
	*s.ptr = generics.PointerTo(int(v))
	return err
}

func (s *intRefValue) Type() string {
	return "*int"
}

func (s *intRefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return strconv.FormatInt(int64(**s.ptr), 10)
}

func intRefConv(sval *intRefValue) (*int, error) {
	return *sval.ptr, nil
}

// GetIntRef return the int ref value of a flag with the given name
func GetIntRef(f *pflag.FlagSet, name string) (*int, error) {

	val, err := getFlagType(f, name, "*int", intRefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// IntRefVar defines a int ref flag with specified name, default value, and usage int.
// The argument p points to a int pointer variable in which to store the value of the flag as referenced value.
func IntRefVar(f *pflag.FlagSet, p **int, name string, value *int, usage string) {
	f.VarP(newIntRefValue(value, p), name, "", usage)
}

// IntRefVarP is like IntRefVar, but accepts a shorthand letter that can be used after a single dash.
func IntRefVarP(f *pflag.FlagSet, p **int, name, shorthand string, value *int, usage string) {
	f.VarP(newIntRefValue(value, p), name, shorthand, usage)
}

// IntRef defines a int ref flag with specified name, default value, and usage int.
// The return value is the address of a int pointer variable that stores the value of the flag.
func IntRef(f *pflag.FlagSet, name string, value *int, usage string) **int {
	p := new(*int)
	IntRefVarP(f, p, name, "", value, usage)
	return p
}

// IntRefP is like IntRef, but accepts a shorthand letter that can be used after a single dash.
func IntRefP(f *pflag.FlagSet, name, shorthand string, value *int, usage string) **int {
	p := new(*int)
	IntRefVarP(f, p, name, shorthand, value, usage)
	return p
}

// IntRefVarPF is like IntRefVarP, but returns the created flag.
func IntRefVarPF(f *pflag.FlagSet, p **int, name, shorthand string, value *int, usage string) *pflag.Flag {
	return f.VarPF(newIntRefValue(value, p), name, shorthand, usage)
}

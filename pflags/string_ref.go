package pflags

import (
	"github.com/spf13/pflag"
)

// -- string Ref Value
type stringRefValue refValue[string]

func newStringRefValue(val *string, p **string) *stringRefValue {
	*p = val
	return &stringRefValue{ptr: p}
}

func (s *stringRefValue) Set(val string) error {
	*s.ptr = &val
	return nil
}
func (s *stringRefValue) Type() string {
	return "*string"
}

func (s *stringRefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return **s.ptr
}

func stringRefConv(sval *stringRefValue) (*string, error) {
	return *sval.ptr, nil
}

// GetStringRef return the string ref value of a flag with the given name
func GetStringRef(f *pflag.FlagSet, name string) (*string, error) {

	val, err := getFlagType(f, name, "*string", stringRefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// StringRefVar defines a string ref flag with specified name, default value, and usage string.
// The argument p points to a string pointer variable in which to store the value of the flag as referenced value.
func StringRefVar(f *pflag.FlagSet, p **string, name string, value *string, usage string) {
	f.VarP(newStringRefValue(value, p), name, "", usage)
}

// StringRefVarP is like StringRefVar, but accepts a shorthand letter that can be used after a single dash.
func StringRefVarP(f *pflag.FlagSet, p **string, name, shorthand string, value *string, usage string) {
	f.VarP(newStringRefValue(value, p), name, shorthand, usage)
}

// StringRef defines a string ref flag with specified name, default value, and usage string.
// The return value is the address of a string pointer variable that stores the value of the flag.
func StringRef(f *pflag.FlagSet, name string, value *string, usage string) **string {
	p := new(*string)
	StringRefVarP(f, p, name, "", value, usage)
	return p
}

// StringRefP is like StringRef, but accepts a shorthand letter that can be used after a single dash.
func StringRefP(f *pflag.FlagSet, name, shorthand string, value *string, usage string) **string {
	p := new(*string)
	StringRefVarP(f, p, name, shorthand, value, usage)
	return p
}

// StringRefVarPF is like StringRefVarP, but returns the created flag.
func StringRefVarPF[T ~map[string]string](f *pflag.FlagSet, p **string, name, shorthand string, value *string, usage string) *pflag.Flag {
	return f.VarPF(newStringRefValue(value, p), name, shorthand, usage)
}

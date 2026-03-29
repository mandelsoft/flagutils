package pflags

import (
	"time"

	"github.com/spf13/pflag"
)

// -- Duration Ref Value
type DurationRefValue refValue[time.Duration]

func newDurationRefValue(val *time.Duration, p **time.Duration) *DurationRefValue {
	*p = val
	return &DurationRefValue{ptr: p}
}

func (s *DurationRefValue) Set(val string) error {
	v, err := time.ParseDuration(val)
	*s.ptr = &v
	return err
}

func (s *DurationRefValue) Type() string {
	return "*Duration"
}

func (s *DurationRefValue) String() string {
	if *s.ptr == nil {
		return "<none>"
	}
	return (**s.ptr).String()
}

func DurationRefConv(sval *DurationRefValue) (*time.Duration, error) {
	return *sval.ptr, nil
}

// GetDurationRef return the time.Duration ref value of a flag with the given name
func GetDurationRef(f *pflag.FlagSet, name string) (*time.Duration, error) {

	val, err := getFlagType(f, name, "*Duration", DurationRefConv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// DurationRefVar defines a time.Duration ref flag with specified name, default value, and usage Duration.
// The argument p points to a Duration pointer variable in which to store the value of the flag as referenced value.
func DurationRefVar(f *pflag.FlagSet, p **time.Duration, name string, value *time.Duration, usage string) {
	f.VarP(newDurationRefValue(value, p), name, "", usage)
}

// DurationRefVarP is like DurationRefVar, but accepts a shorthand letter that can be used after a single dash.
func DurationRefVarP(f *pflag.FlagSet, p **time.Duration, name, shorthand string, value *time.Duration, usage string) {
	f.VarP(newDurationRefValue(value, p), name, shorthand, usage)
}

// DurationRef defines a time.Duration ref flag with specified name, default value, and usage Duration.
// The return value is the address of a Duration pointer variable that stores the value of the flag.
func DurationRef(f *pflag.FlagSet, name string, value *time.Duration, usage string) **time.Duration {
	p := new(*time.Duration)
	DurationRefVarP(f, p, name, "", value, usage)
	return p
}

// DurationRefP is like DurationRef, but accepts a shorthand letter that can be used after a single dash.
func DurationRefP(f *pflag.FlagSet, name, shorthand string, value *time.Duration, usage string) **time.Duration {
	p := new(*time.Duration)
	DurationRefVarP(f, p, name, shorthand, value, usage)
	return p
}

// DurationRefVarPF is like DurationRefVarP, but returns the created flag.
func DurationRefVarPF(f *pflag.FlagSet, p **time.Duration, name, shorthand string, value *time.Duration, usage string) *pflag.Flag {
	return f.VarPF(newDurationRefValue(value, p), name, shorthand, usage)
}

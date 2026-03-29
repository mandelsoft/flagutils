package pflags

import (
	"fmt"

	"github.com/mandelsoft/goutils/generics"
	"github.com/spf13/pflag"
)

// func to return a given type for a given flag name
func getFlagType[T, R any](f *pflag.FlagSet, name, ftype string, convFunc func(T) (R, error)) (R, error) {
	var _nil R
	flag := f.Lookup(name)
	if flag == nil {
		return _nil, f.MarkHidden(name) // enforce creation of correct (private) error
	}

	if v, ok := flag.Value.(T); !ok {
		err := fmt.Errorf("trying to get %s value of flag of type %s", ftype, flag.Value.Type())
		return _nil, err
	} else {
		if convFunc != nil {
			return convFunc(v)
		}
		return generics.TryCastE[R](flag.Value)
	}
}

// -- Ref Value
type refValue[T any] struct {
	ptr **T
}

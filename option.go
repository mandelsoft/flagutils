package flagutils

import (
	"context"
	"github.com/mandelsoft/goutils/set"
	"github.com/spf13/pflag"
)

// Options provides an interface for adding flags to a given pflag.FlagSet
// using the AddFlags method.
type Options interface {
	AddFlags(fs *pflag.FlagSet)
}

// Validation defines an interface for objects that can be validated based on
// an OptionSet  within a given context.
// Optionally, the given context as well as the other options in the OptionSet
// can also be used to complete the option state.
// If nested elements are used, they must be validated using the given
// ValidationSet to assert they are already validated before used.
type Validation interface {
	Validate(ctx context.Context, opts OptionSet, v ValidationSet) error
}

// Usage is an interface representing an entity capable of producing a usage
// string via the Usage method.
type Usage interface {
	Usage() string
}

// ValidationSet is a set of Validation elements that ensures each element
// is validated only once within a context. It keeps a set of already
// validated objects. If there are cyclic evaluations, only the first call
// evaluates the object. The order therefore depends on the order of the
// executed initial validations, No error is provided for such cyclic scenarios.
type ValidationSet set.Set[Validation]

func (s ValidationSet) Validate(ctx context.Context, opts OptionSet, o any) error {
	if v, ok := o.(Validation); ok {
		if !set.Set[Validation](s).Has(v) {
			set.Set[Validation](s).Add(v)
			return v.Validate(ctx, opts, s)
		}
	}
	return nil
}

func ValidatedOptions[O any](ctx context.Context, opts OptionSet, s ValidationSet) (O, error) {
	var _nil O
	o := GetFrom[O](opts)

	var a any = o
	if a != nil {
		err := s.Validate(ctx, opts, o)
		if err != nil {
			return _nil, err
		}
	}
	return o, nil
}

type OptionBase[T Options] struct {
	self  T
	long  *string
	short *string
}

func NewBase[T Options](self T) OptionBase[T] {
	return OptionBase[T]{self: self}
}

func (o *OptionBase[T]) Long(def string) string {
	if o.long == nil {
		return def
	}
	return *o.long
}

func (o *OptionBase[T]) Short(def string) string {
	if o.short == nil {
		return def
	}
	return *o.short
}

func (o *OptionBase[T]) WithNames(l, s string) T {
	o.long = &l
	o.short = &s
	return o.self
}

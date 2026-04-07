package flagutils

import (
	"context"

	"github.com/mandelsoft/goutils/reflectutils"
	"github.com/mandelsoft/goutils/set"
	"github.com/modern-go/reflect2"
)

// Validatable defines an interface for objects that can be validated based on
// an OptionSet  within a given context.
// Optionally, the given context as well as the other options in the OptionSet
// can also be used to complete the option state.
// If nested elements are used, they must be validated using the given
// ValidationSet to assert they are already validated before used.
type Validatable interface {
	Validate(ctx context.Context, opts OptionSet, v ValidationSet) error
}

// ValidationSet is a set of Validatable elements that ensures each element
// is validated only once within a context. It keeps a set of already
// validated objects. If there are cyclic evaluations, only the first call
// evaluates the object. The order therefore depends on the order of the
// executed initial validations, No error is provided for such cyclic scenarios.
type ValidationSet set.Set[Validatable]

func (s ValidationSet) Validate(ctx context.Context, opts OptionSet, orig any) error {
	o := orig
	for o != nil {
		if v, ok := o.(Validatable); ok {
			if !set.Set[Validatable](s).Has(v) {
				set.Set[Validatable](s).Add(v)
				return v.Validate(ctx, opts, s)
			}
			return nil
		}
		o = reflectutils.UnwrapAny(o)
	}

	o = orig
	for o != nil {
		if v, ok := o.(OptionSetProvider); ok {
			return s.ValidateSet(ctx, opts, v)
		}
		o = reflectutils.UnwrapAny(o)
	}
	return nil
}

// ValidateSet validates the OptionSet given by an OptionSetProvider against a more
// general OptionSet using the provided ValidationSet.
// It iterates over the options in the set and applies validation using the
// provided context and the general OptionSet.
// If validation fails for any option, the function returns the respective error.
// This function is intended to be used by Validate method in some Options
// object requiring to forward validation to a nested OptionSet.
// Note: If an object implements a Validate method, it is also responsible to handle nested options.
func (s ValidationSet) ValidateSet(ctx context.Context, opts OptionSet, set OptionSetProvider) error {
	if v, ok := set.(Validatable); ok {
		err := v.Validate(ctx, opts, s)
		if err != nil {
			return err
		}
	} else {
		for o := range set.AsOptionSet().Options {
			err := s.Validate(ctx, opts, o)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidatedOptions provides a validated Options object of the given type.
// The type is typically a pointer type to the Options struct.
func ValidatedOptions[O any](ctx context.Context, opts OptionSet, s ValidationSet) (O, error) {
	var _nil O
	o := GetFrom[O](opts)

	if !reflect2.IsNil(o) {
		err := s.Validate(ctx, opts, o)
		if err != nil {
			return _nil, err
		}
	}
	return o, nil
}

// ValidatedFilteredOptions filters an OptionSet for elements of type O and validates each element against a ValidationSet.
// Returns the filtered and validated list of elements or an error if validation fails.
// It is typically used with an interface type to get all Options objects implementing this interface.
func ValidatedFilteredOptions[O any](ctx context.Context, opts OptionSet, s ValidationSet) ([]O, error) {
	list := Filter[O](opts)

	for _, o := range list {
		err := s.Validate(ctx, opts, o)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

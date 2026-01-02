package flagutils

import (
	"context"
	"github.com/mandelsoft/goutils/iterutils"
	"github.com/mandelsoft/goutils/set"
	"github.com/modern-go/reflect2"
	"github.com/spf13/pflag"
)

// Options provides an interface for adding flags to a given pflag.FlagSet
// using the AddFlags method.
type Options interface {
	AddFlags(fs *pflag.FlagSet)
}

// Usage is an interface representing an entity capable of producing a usage
// string via the Usage method. This info is a length description of the
// purpose of the option, which can be used in the command description.
type Usage interface {
	Usage() string
}

////////////////////////////////////////////////////////////////////////////////

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

func (s ValidationSet) Validate(ctx context.Context, opts OptionSet, o any) error {
	if v, ok := o.(Validatable); ok {
		if !set.Set[Validatable](s).Has(v) {
			set.Set[Validatable](s).Add(v)
			return v.Validate(ctx, opts, s)
		}
	}
	return nil
}

// ValidateSet validates the OptionSet givey an OptionSetProvider against a more
// general OptionSet using the provided ValidationSet.
// It iterates over the options in the set and applies validation using the
// provided context and the general OptionSet.
// If validation fails for any option, the function returns the respective error.
// This function is intended to be used by Validation method in some Options
// object requiring to forward Validation to a nested OptionSet.
// Note: If an object implements a Validation method, it is also responsible to handle nested options.
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

////////////////////////////////////////////////////////////////////////////////

// Finalizable represents a type that can perform a finalization operation with
// a context and a set of options. Options keeping external state should implement
// this interface.
type Finalizable interface {
	Finalize(ctx context.Context, opts OptionSet, v FinalizationSet) error
}

// FinalizationSet is a set of finalization elements that ensures each element
// is finalized only once within a context. It keeps a set of already
// finalized objects. If there are cyclic finalizations, only the first call
// finalizes the object. The order therefore depends on the order of the
// executed initial finalizations, No error is provided for such cyclic scenarios.
type FinalizationSet set.Set[Finalizable]

func (s FinalizationSet) Finalize(ctx context.Context, opts OptionSet, o any) error {
	if v, ok := o.(Finalizable); ok {
		if !set.Set[Finalizable](s).Has(v) {
			set.Set[Finalizable](s).Add(v)
			return v.Finalize(ctx, opts, s)
		}
	}
	return nil
}

// FinalizeSet finalizes the OptionSet given by the OptionSetProvider against a
// more general OptionSet using the provided FinalizationSet.
// It iterates over the options in the set and applies validation using the
// provided context and the general OptionSet.
// If validation fails for any option, the function returns the respective error.
// This function is intended to be used by Validation method in some Options
// object requiring to forward Validation to a nested OptionSet.
// Note: If an object implements a Finalize method, it is also responsible to handle nested options.
func (s FinalizationSet) FinalizeSet(ctx context.Context, opts OptionSet, set OptionSetProvider) error {
	if v, ok := set.(Finalizable); ok {
		err := v.Finalize(ctx, opts, s)
		if err != nil {
			return err
		}
	} else {
		for o := range iterutils.Reverse(set.AsOptionSet().Options) {
			err := s.Finalize(ctx, opts, o)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

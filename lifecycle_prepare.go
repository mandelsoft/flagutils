package flagutils

import (
	"context"

	"github.com/mandelsoft/goutils/reflectutils"
	"github.com/mandelsoft/goutils/set"
	"github.com/modern-go/reflect2"
)

// Preparable defines an interface for objects that can be prepared based on
// an OptionSet  within a given context.
// Optionally, the given context as well as the other options in the OptionSet
// can also be used to complete the option state.
// If nested elements are used, they must be prepared using the given
// PreparationSet to assert they are already prepared before used.
type Preparable interface {
	Prepare(ctx context.Context, opts OptionSet, v PreparationSet) error
}

// PreparationSet is a set of Preparable elements that ensures each element
// is prepared only once within a context. It keeps a set of already
// prepared objects. If there are cyclic evaluations, only the first call
// evaluates the object. The order therefore depends on the order of the
// executed initial preparations, No error is provided for such cyclic scenarios.
type PreparationSet set.Set[Preparable]

func (s PreparationSet) Prepare(ctx context.Context, opts OptionSet, orig any) error {
	o := orig
	for o != nil {
		if v, ok := o.(Preparable); ok {
			if !set.Set[Preparable](s).Has(v) {
				set.Set[Preparable](s).Add(v)
				return v.Prepare(ctx, opts, s)
			}
			return nil
		}
		o = reflectutils.UnwrapAny(o)
	}

	o = orig
	for o != nil {
		if v, ok := o.(OptionSetProvider); ok {
			return s.PrepareSet(ctx, opts, v)
		}
		o = reflectutils.UnwrapAny(o)
	}
	return nil
}

// PrepareSet validates the OptionSet given by an OptionSetProvider against a more
// general OptionSet using the provided PreparationSet.
// It iterates over the options in the set and applies preparation using the
// provided context and the general OptionSet.
// If preparation fails for any option, the function returns the respective error.
// This function is intended to be used by Prepare methods in some Options
// object requiring to forward preparation to a nested OptionSet.
// Note: If an object implements a Prepare method, it is also responsible to handle nested options.
func (s PreparationSet) PrepareSet(ctx context.Context, opts OptionSet, set OptionSetProvider) error {
	if v, ok := set.(Preparable); ok {
		err := v.Prepare(ctx, opts, s)
		if err != nil {
			return err
		}
	} else {
		for o := range set.AsOptionSet().Options {
			err := s.Prepare(ctx, opts, o)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// PreparedOptions provides a validated Options object of the given type.
// The type is typically a pointer type to the Options struct.
func PreparedOptions[O any](ctx context.Context, opts OptionSet, s PreparationSet) (O, error) {
	var _nil O
	o := GetFrom[O](opts)

	if !reflect2.IsNil(o) {
		err := s.Prepare(ctx, opts, o)
		if err != nil {
			return _nil, err
		}
	}
	return o, nil
}

// PreparedFilteredOptions filters an OptionSet for elements of type O and prepares each element against a PreparationSet.
// Returns the filtered and prepared list of elements or an error if preparation fails.
// It is typically used with an interface type to get all Options objects implementing this interface.
func PreparedFilteredOptions[O any](ctx context.Context, opts OptionSet, s PreparationSet) ([]O, error) {
	list := Filter[O](opts)

	for _, o := range list {
		err := s.Prepare(ctx, opts, o)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

package flagutils

import (
	"context"

	"github.com/mandelsoft/goutils/iterutils"
	"github.com/mandelsoft/goutils/set"
	"github.com/modern-go/reflect2"
)

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

// FinalizedOptions provides a finalized Options object of the given type.
// The type is typically a pointer type to the Options struct.
func FinalizedOptions[O any](ctx context.Context, opts OptionSet, s FinalizationSet) (O, error) {
	var _nil O
	o := GetFrom[O](opts)

	if !reflect2.IsNil(o) {
		err := s.Finalize(ctx, opts, o)
		if err != nil {
			return _nil, err
		}
	}
	return o, nil
}

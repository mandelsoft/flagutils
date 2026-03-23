package flagutils

import "context"

// Prepare checks whether the provided OptionSetProvider or its nested options
// implement the Preparable interface and prepares them.
// It returns an error if any preparation fails or nil if all preparations succeed.
// Preparation should be called after the setup of an OptionSet and before
// it is added to a pflag.FlagSet. It is intended to link various options,
// check whether thay are compatible and prepare dependent default values
// and/or values helps.
func Prepare(ctx context.Context, set OptionSetProvider, val PreparationSet) error {
	if val == nil {
		val = PreparationSet{}
	}
	base := set.AsOptionSet()
	if v, ok := set.(Preparable); ok {
		err := v.Prepare(ctx, base, val)
		if err != nil {
			return err
		}
	} else {
		return val.PrepareSet(ctx, base, base)
	}
	return nil
}

// Validate checks whether the provided OptionSetProvider or its nested options
// implement the Validatable interface and validates them.
// It returns an error if any validation fails or nil if all validations succeed.
func Validate(ctx context.Context, set OptionSetProvider, val ValidationSet) error {
	if val == nil {
		val = ValidationSet{}
	}
	base := set.AsOptionSet()
	if v, ok := set.(Validatable); ok {
		err := v.Validate(ctx, base, val)
		if err != nil {
			return err
		}
	} else {
		return val.ValidateSet(ctx, base, base)
	}
	return nil
}

// Finalize checks whether the provided OptionSetProvider or its nested options
// implement the Finalizable interface and finalizes them.
// It returns an error if any finalization fails or nil if all finalizations succeed.
func Finalize(ctx context.Context, set OptionSetProvider, val FinalizationSet) error {
	if val == nil {
		val = FinalizationSet{}
	}
	return val.FinalizeSet(ctx, set.AsOptionSet(), set)
}

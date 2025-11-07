package flagutils

import (
	"context"
	"github.com/spf13/pflag"
	"reflect"
)

// OptionSetProvider defines an interface for types capable of providing an OptionSet by implementing the AsOptionSet method.
type OptionSetProvider interface {
	AsOptionSet() OptionSet
}

// OptionSet is an interface representing a set of Options. It acts as Options and OptionSetProvider
// and provides a method to iterate over nested Options.
type OptionSet interface {
	Options
	OptionSetProvider
	Options(yield func(Options) bool)
}

// DefaultOptionSet defines a slice of Options, representing a basic
// implementation of an OptionSet.
type DefaultOptionSet []Options

func (s DefaultOptionSet) AsOptionSet() OptionSet {
	return s
}

var _ OptionSet = (*DefaultOptionSet)(nil)

func (s *DefaultOptionSet) Add(o ...Options) *DefaultOptionSet {
	*s = append(*s, o...)
	return s
}

func (s DefaultOptionSet) AddFlags(fs *pflag.FlagSet) {
	for _, o := range s {
		o.AddFlags(fs)
	}
}

func (s DefaultOptionSet) Options(yield func(Options) bool) {
	for _, o := range s {
		if !yield(o) {
			return
		}
	}
}

func (s DefaultOptionSet) Usage() string {
	u := ""
	for _, n := range s {
		if c, ok := n.(Usage); ok {
			u += c.Usage()
		}
	}
	return u
}

func get(pv reflect.Value, o any) bool {
	ov := reflect.ValueOf(o)
	if pv.Elem().Kind() == reflect.Ptr {
		// pointer to pointer
		if ov.Type() == pv.Elem().Type() {
			pv.Elem().Set(ov)
			return true
		}
	} else {
		// pointer to value
		if ov.Type().AssignableTo(pv.Type()) {
			pv.Elem().Set(ov.Elem())
			return true
		}
		if ov.Type().AssignableTo(pv.Type().Elem()) {
			pv.Elem().Set(ov)
			return true
		}
	}
	return false
}

func retrieveFrom(set OptionSetProvider, pv reflect.Value) bool {
	if get(pv, set) {
		return true
	}
	for o := range set.AsOptionSet().Options {
		if set, ok := o.(OptionSetProvider); ok {
			if ok := retrieveFrom(set, pv); ok {
				return true
			}
		} else {
			if get(pv, o) {
				return true
			}
		}
	}
	return false
}

// RetrieveFrom extracts the option for a given target. This might be a
//   - pointer to a struct implementing the Options interface which
//     will fill the struct with a copy of the options OR
//   - a pointer to such a pointer which will be filled with the
//     pointer to the actual member of the OptionSet.
func RetrieveFrom(set OptionSetProvider, proto interface{}) bool {
	return retrieveFrom(set, reflect.ValueOf(proto))
}

// GetFrom retrieves an option of type T from the provided OptionSetProvider
// and returns it. T is typically a pointer to an option struct of type Options.
// If an interface type is used the first found implementation is returned.
// To get all options implementing an interface use Filter.
func GetFrom[T any](set OptionSetProvider) T {
	var r T
	RetrieveFrom(set, &r)
	return r
}

func filter[T any](set OptionSetProvider, result *[]T) {
	if o, ok := set.(Options); ok {
		if v, ok := o.(T); ok {
			*result = append(*result, v)
		}
	}
	for o := range set.AsOptionSet().Options {
		if v, ok := o.(OptionSetProvider); ok {
			filter[T](v, result)
		} else {
			if v, ok := o.(T); ok {
				*result = append(*result, v)
			}
		}
	}
}

// Filter extracts elements of type T from the provided OptionSetProvider and returns them as a slice of T.
func Filter[T any](set OptionSetProvider) []T {
	var result []T
	filter[T](set, &result)
	return result
}

// Validate checks whether the provided OptionSetProvider or its nested options
// implement the Validation interface and validates them.
// It returns an error if any validation fails or nil if all validations succeed.
func Validate(ctx context.Context, set OptionSetProvider, val ValidationSet) error {
	if val == nil {
		val = ValidationSet{}
	}
	base := set.AsOptionSet()
	if v, ok := set.(Validation); ok {
		err := v.Validate(ctx, base, val)
		if err != nil {
			return err
		}
	} else {
		for o := range set.AsOptionSet().Options {
			if v, ok := o.(Validation); ok {
				err := v.Validate(ctx, base, val)
				if err != nil {
					return err
				}
			}
		}
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
	base := set.AsOptionSet()
	if v, ok := set.(Finalizable); ok {
		err := v.Finalize(ctx, base, val)
		if err != nil {
			return err
		}
	} else {
		for o := range set.AsOptionSet().Options {
			if v, ok := o.(Finalizable); ok {
				err := v.Finalize(ctx, base, val)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

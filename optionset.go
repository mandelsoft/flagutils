package flagutils

import (
	"reflect"
)

// OptionSetProvider defines an interface for types capable of providing an OptionSet by implementing the AsOptionSet method.
type OptionSetProvider interface {
	AsOptionSet() OptionSet
}

// OptionSet is an interface representing a set of Options. It acts as Options and OptionSetProvider
// and provides a method to iterate over nested Options.
// There is an intended lifecycle for an OptionSet:
//   - Composition of the set
//   - Preparation using Prepare and a PreparationSet to complete the option definition incorporation other
//     options of the final OptionSet.
//   - Apply to pflag.FlagSet
//   - Evaluate on current command line options.
//   - Validation using Validate and a ValidationSet to validate the settings and prepare some state usable by
//     the intended application.
//   - (Run the application using the options (potetially with the From calls from various options to retrieve
//     them from the OptionSet.
//   - Finalization using Finalize and a FinalizationSet to cleanup temporary state build during Validation.
type OptionSet interface {
	Options
	OptionSetProvider
	Options(yield func(Options) bool)
}

type ExtendableOptionSet interface {
	OptionSet
	Add(o ...Options) ExtendableOptionSet
}

////////////////////////////////////////////////////////////////////////////////

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

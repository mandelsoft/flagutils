package flagutils

import (
	"fmt"
	"reflect"

	"github.com/mandelsoft/goutils/matcher"
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
//   - (Run the application using the options (potentially with the From calls from various options to retrieve
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

func GetFrom2[T any](set OptionSetProvider) (T, bool) {
	var r T
	ok := RetrieveFrom(set, &r)
	return r, ok
}

// GetFilteredFrom is like GetFrom, but uses an instance filter
// to select an instance for a type supporting different instance specific
// flag names.
func GetFilteredFrom[T any](set OptionSetProvider, filter matcher.Matcher[T]) T {
	var _nil T
	list := Filter[T](set, filter)
	if len(list) == 0 {
		return _nil
	}
	return list[0]
}

func GetFilteredFrom2[T any](set OptionSetProvider, filter matcher.Matcher[T]) (T, bool) {
	var _nil T
	list := Filter[T](set, filter)
	if len(list) == 0 {
		return _nil, false
	}
	return list[0], true
}

func filter[T any](set OptionSetProvider, result *[]T, check matcher.Matcher[T]) {
	if o, ok := set.(Options); ok {
		if v, ok := o.(T); ok {
			*result = append(*result, v)
		}
	}
	for o := range set.AsOptionSet().Options {
		if v, ok := o.(OptionSetProvider); ok {
			filter[T](v, result, check)
		} else {
			if v, ok := o.(T); ok && check(v) {
				*result = append(*result, v)
			}
		}
	}
}

// Filter extracts elements of type T from the provided OptionSetProvider and returns them as a slice of T.
// An optional matcher can be used to additionally filter the result.
func Filter[T any](set OptionSetProvider, check ...matcher.Matcher[T]) []T {
	var result []T
	filter[T](set, &result, matcher.And(check...))
	return result
}

// Assure expects an ExtendableOptionSet and adds an options object
// provided by a factory method, if it is not yet present.
// Optional matchers can be used to apply additional filters,
// for example an instance filter for types applicable for different option names.
func Assure[T Options](opts OptionSet, f func() T, check ...matcher.Matcher[T]) error {
	list := Filter[T](opts, check...)
	if len(list) > 0 {
		return nil
	}
	if m, ok := opts.(ExtendableOptionSet); ok {
		m.Add(f())
		return nil
	}
	return fmt.Errorf("option set must implement flagutils.ExtendableOptionSet")
}

func SetAssured[T Options](tgt *T, opts OptionSet, f func() T, check ...matcher.Matcher[T]) error {
	list := Filter[T](opts, check...)
	if len(list) > 0 {
		*tgt = list[0]
		return nil
	}
	if m, ok := opts.(ExtendableOptionSet); ok {
		*tgt = f()
		m.Add(*tgt)
		return nil
	}
	return fmt.Errorf("option set must implement flagutils.ExtendableOptionSet")
}

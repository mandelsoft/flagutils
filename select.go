package flagutils

import "github.com/mandelsoft/goutils/sliceutils"

type OptionSelector func(Options) bool

func Not(s OptionSelector) OptionSelector {
	return func(o Options) bool {
		return !s(o)
	}
}

func Or(s ...OptionSelector) OptionSelector {
	return func(o Options) bool {
		for _, sel := range s {
			if sel(o) {
				return true
			}
		}
		return false
	}
}

func And(s ...OptionSelector) OptionSelector {
	return func(o Options) bool {
		for _, sel := range s {
			if !sel(o) {
				return false
			}
		}
		return true
	}
}

func Always() OptionSelector {
	return func(o Options) bool {
		return true
	}
}

func Never() OptionSelector {
	return func(o Options) bool {
		return false
	}
}

////////////////////////////////////////////////////////////////////////////////

func Implements[T any](o Options) bool {
	_, ok := o.(T)
	return ok
}

////////////////////////////////////////////////////////////////////////////////

func Select(set OptionSet, sel OptionSelector) DefaultOptionSet {
	var result DefaultOptionSet

	if o, ok := set.(Options); ok {
		if sel(o) {
			result = append(result, o)
		}
	}
	for o := range set.Options {
		if s, ok := o.(OptionSet); ok {
			result = append(result, Select(s, sel))
		} else {
			if sel(o) {
				result = append(result, o)
			}
		}
	}
	return result
}

func SelectByInterface[T any](set OptionSet, sel ...OptionSelector) []T {
	r := Select(set, And(Implements[T], And(sel...)))
	return sliceutils.Convert[T](r)
}

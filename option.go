package flagutils

import (
	"context"

	"github.com/mandelsoft/goutils/matcher"
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

type OptionsRef[T Options] struct {
	factory func() T
	matcher matcher.Matcher[T]
	Options T
}

var (
	_ Options    = (*OptionsRef[Options])(nil)
	_ Preparable = (*OptionsRef[Options])(nil)
)

// NewOptionsRef creates a new dynamic reference to Options
// added once to an OptionSet on the fly during the preparation phase.
// This can be used to handled Options shared among other
// aggregating Options structures.
func NewOptionsRef[T Options](f func() T, check ...matcher.Matcher[T]) *OptionsRef[T] {
	return &OptionsRef[T]{factory: f, matcher: matcher.And(check...)}
}

func (o *OptionsRef[T]) AddFlags(fs *pflag.FlagSet) {
	// just a marker method
	// Option settings are deferred ro assured options.
}

func (o *OptionsRef[T]) Prepare(ctx context.Context, opts OptionSet, v PreparationSet) error {
	return SetAssured(&o.Options, opts, o.factory, o.matcher)
}

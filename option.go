package flagutils

import (
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

////////////////////////////////////////////////////////////////////////////////

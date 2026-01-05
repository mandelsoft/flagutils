package flagutils

import "github.com/spf13/pflag"

type NoOptions struct{}

var _ Options = (*NoOptions)(nil)

func (NoOptions) AddFlags(fs *pflag.FlagSet) {}

////////////////////////////////////////////////////////////////////////////////

type SetBasedOptions struct {
	set DefaultOptionSet
}

func (s *SetBasedOptions) AddFlags(fs *pflag.FlagSet) {
	s.set.AddFlags(fs)
}

func (s *SetBasedOptions) Add(o ...Options) OptionSet {
	s.set.Add(o...)
	return s.set
}

func (s *SetBasedOptions) AsOptionSet() OptionSet {
	return s.set
}

////////////////////////////////////////////////////////////////////////////////

// DefaultOptionSet defines a slice of Options, representing a basic
// implementation of an OptionSet.
type DefaultOptionSet []Options

func (s DefaultOptionSet) AsOptionSet() OptionSet {
	return s
}

var _ ExtendableOptionSet = (*DefaultOptionSet)(nil)

func (s *DefaultOptionSet) Add(o ...Options) ExtendableOptionSet {
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

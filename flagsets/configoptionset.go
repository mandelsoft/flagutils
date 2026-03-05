package flagsets

import (
	"fmt"
	"sync"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/maputils"
)

// OptionType describes a particular type of option
// It has a name and a description, and can be
// used to create an Option instance with the same name, which can
// be added to a pflag.FlagSet,
// Every Optiontype has a technical type (the type of the
// underlying pflag.Flag.
// Two typed are identical if name, description
// and technical type are identical.
type OptionType interface {
	GetName() string
	GetDescription() string

	Create() Option

	Equal(optionType OptionType) bool
}

// OptionTypeSet represents the type
// for a set of [Option]s by describing
// the set of [OptionType]s.
// It has a name and nested [OptionType]s.
// This nesting could be described by
// other [OptionTypeSet]s.
// [OptionType]s hereby, might be shared with nested
// Sets, if their technical type matches.
// In an OptionTypeSet the names must be unique.
// If a nested set contains the same name as
// the nesting one (or two nested sets contain the
// same name) the option types must be identical.
type OptionTypeSet interface {
	AddGroups(groups ...string)

	GetName() string

	Size() int
	OptionTypes() []OptionType
	OptionTypeNames() []string
	SharedOptionTypes() []OptionType

	HasOptionType(name string) bool
	HasSharedOptionType(name string) bool

	GetSharedOptionType(name string) OptionType
	GetOptionType(name string) OptionType
	GetTypeSet(name string) OptionTypeSet
	OptionTypeSets() []OptionTypeSet

	AddOptionType(OptionType) error
	AddTypeSet(OptionTypeSet) error
	AddAll(o OptionTypeSet) (duplicated OptionTypeSet, err error)

	Close(funcs ...func([]OptionType) error) error

	CreateOptions() Options
	AddGroupsToOption(o Option)
}

type ptionTypeSet struct {
	lock    sync.RWMutex
	name    string
	options map[string]OptionType
	sets    map[string]OptionTypeSet
	shared  map[string][]OptionTypeSet
	groups  []string

	closed bool
}

func NewOptionTypeSet(name string, types ...OptionType) OptionTypeSet {
	set := &ptionTypeSet{
		name:    name,
		options: map[string]OptionType{},
		sets:    map[string]OptionTypeSet{},
		shared:  map[string][]OptionTypeSet{},
	}
	for _, t := range types {
		set.AddOptionType(t)
	}
	return set
}

func (s *ptionTypeSet) AddGroups(groups ...string) {
	s.groups = AddGroups(s.groups, groups...)
}

func (s *ptionTypeSet) Close(funcs ...func([]OptionType) error) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(funcs) > 0 {
		list := s.optionTypes()
		for _, f := range funcs {
			if f != nil {
				err := f(list)
				if err != nil {
					return err
				}
			}
		}
	}
	s.closed = true
	return nil
}

func (s *ptionTypeSet) GetName() string {
	return s.name
}

func (s *ptionTypeSet) AddOptionType(optionType OptionType) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.closed {
		return errors.ErrClosed("config option set")
	}
	name := optionType.GetName()
	s.options[name] = optionType
	return nil
}

func (s *ptionTypeSet) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.options)
}

func (s *ptionTypeSet) OptionTypes() []OptionType {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.optionTypes()
}

func (s *ptionTypeSet) optionTypes() []OptionType {
	var list []OptionType
	for _, o := range s.options {
		list = append(list, o)
	}
	return list
}

func (s *ptionTypeSet) OptionTypeNames() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return maputils.OrderedKeys(s.options)
}

func (s *ptionTypeSet) SharedOptionTypes() []OptionType {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var list []OptionType
	for n, o := range s.options {
		if _, ok := s.shared[n]; ok {
			list = append(list, o)
		}
	}
	return list
}

func (s *ptionTypeSet) HasOptionType(name string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, ok := s.options[name]
	return ok
}

func (s *ptionTypeSet) HasSharedOptionType(name string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, ok := s.shared[name]
	return ok
}

func (s *ptionTypeSet) GetOptionType(name string) OptionType {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.options[name]
}

func (s *ptionTypeSet) GetSharedOptionType(name string) OptionType {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if _, ok := s.shared[name]; ok {
		return s.options[name]
	}
	return nil
}

func (s *ptionTypeSet) AddTypeSet(set OptionTypeSet) error {
	if set == nil {
		return nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.closed {
		return errors.ErrClosed("config option set")
	}

	name := set.GetName()
	if nested, ok := s.sets[name]; ok {
		if nested == set {
			return nil
		}
		return fmt.Errorf("%s: config type set with name %q already added", s.GetName(), name)
	}

	return set.Close(func(list []OptionType) error {
		// check for problem first
		err := s.check(list)
		if err != nil {
			return err
		}
		// now align data structure
		for _, o := range list {
			if _, ok := s.options[o.GetName()]; ok {
				s.shared[o.GetName()] = append(s.shared[o.GetName()], set)
			} else {
				s.options[o.GetName()] = o
			}
		}
		s.sets[name] = set
		return nil
	})
}

func (s *ptionTypeSet) GetTypeSet(name string) OptionTypeSet {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.sets[name]
}

func (s *ptionTypeSet) OptionTypeSets() []OptionTypeSet {
	s.lock.RLock()
	defer s.lock.RUnlock()

	result := make([]OptionTypeSet, 0, len(s.sets))
	for _, t := range s.sets {
		result = append(result, t)
	}
	return result
}

func (s *ptionTypeSet) AddGroupsToOption(o Option) {
	if !s.HasOptionType(o.GetName()) {
		return
	}
	if len(s.groups) > 0 {
		o.AddGroups(s.groups...)
	}
	for _, set := range s.shared[o.GetName()] {
		set.AddGroupsToOption(o)
	}
}

func (s *ptionTypeSet) CreateOptions() Options {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var opts []Option

	for n := range s.options {
		opt := s.options[n].Create()
		s.AddGroupsToOption(opt)
		opts = append(opts, opt)
	}
	return NewOptions(opts)
}

func (s *ptionTypeSet) AddAll(o OptionTypeSet) (duplicates OptionTypeSet, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.closed {
		return nil, errors.ErrClosed("config option set")
	}

	list := o.OptionTypes()
	if err := s.check(list); err != nil {
		return nil, err
	}
	duplicates = NewOptionTypeSet("duplicates")
	for _, t := range list {
		_, ok := s.options[t.GetName()]
		if !ok {
			s.options[t.GetName()] = t
		} else {
			duplicates.AddOptionType(t)
		}
	}
	return duplicates, nil
}

func (s *ptionTypeSet) check(list []OptionType) error {
	for _, o := range list {
		old := s.options[o.GetName()]
		if old != nil && !old.Equal(o) {
			return fmt.Errorf("option type %s doesn't match (%T<->%T)", o.GetName(), o, old)
		}
	}
	return nil
}

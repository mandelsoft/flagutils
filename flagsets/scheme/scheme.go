// Package scheme uses the flagsets.OptionSet handling from package
// [flagsets] to handle a command line option based configuration
// of instances of dynamically registered Type declarations.
package scheme

import (
	"maps"
	"reflect"

	"github.com/mandelsoft/flagutils/flagsets"
	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/generics"
)

type Type[T any] interface {
	flagsets.ConfigOptionTypeSetHandler
	CreateObject() T
}

type _type[T, B any] struct {
	typ reflect.Type
	flagsets.ConfigOptionTypeSetHandler
}

// NewType creates a new Type for an implementation T, which must implement interface B.
// Because of the elaborated type system in Go this cannot be expressed as type constraint.
func NewType[T, B any](name string, adder flagsets.ConfigAdder, types ...flagsets.OptionType) Type[B] {
	typ := generics.TypeOf[T]()
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return &_type[T, B]{
		typ:                        typ,
		ConfigOptionTypeSetHandler: flagsets.NewConfigOptionTypeSetHandler(name, adder, types...),
	}
}

func (t *_type[T, B]) CreateObject() B {
	return reflect.New(t.typ).Interface().(B)
}

////////////////////////////////////////////////////////////////////////////////

// Scheme describes a set of [Type]s and is able
// to create appropriate implementation variants for a flagsets.Config
// object by using json un-/marshalling.
// It uses the flagsets.ConfigOptionTypeSetHandler feature of Type objects
// to contruct ta flagsets.TypedOptionSetConfigProvider with an appropriate
// type options use dto select the variant from the given options.
type Scheme[T any] interface {
	AddType(t Type[T])
	GetTypes() map[string]Type[T]
	CreateOptionSetConfigProvider() (flagsets.TypedOptionSetConfigProvider, error)

	CreateObject(config flagsets.Config) (T, error)
}

func New[T any]() Scheme[T] {
	return _scheme[T]{}
}

type _scheme[T any] map[string]Type[T]

func (s _scheme[T]) AddType(t Type[T]) {
	s[t.GetName()] = t
}

func (s _scheme[T]) GetTypes() map[string]Type[T] {
	return maps.Clone(s)
}

func (s _scheme[T]) CreateOptionSetConfigProvider() (flagsets.TypedOptionSetConfigProvider, error) {
	p := flagsets.NewTypedConfigProvider("object", "some object variants", "objectType")
	for _, v := range s {
		err := p.AddTypeSet(v)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (s _scheme[T]) CreateObject(config flagsets.Config) (T, error) {
	var _nil T
	v := config["type"]
	if v == nil {
		return _nil, errors.New("no type set")
	}
	typ, ok := v.(string)
	if !ok {
		return _nil, errors.New("type must be strint attribute")
	}
	t := s[typ]
	if t == nil {
		return _nil, errors.Newf("unknown type %q", typ)
	}
	e := t.CreateObject()
	err := flagsets.UnmarshalConfig(config, e)
	if err != nil {
		return _nil, err
	}
	return e, nil
}

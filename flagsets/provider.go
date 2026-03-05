package flagsets

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
)

// ConfigProvider is able to create
// a set of command line [Option]s in form
// of an OptionSet and extract a Config from
// Options (might be a super set).
type ConfigProvider interface {
	CreateOptions() Options
	GetConfigFor(opts Options) (Config, error)
}

// TypedOptionSetConfigProvider is ConfigProvider
// based on an OptionTypeSet with a dedicated OptionType
// used to ???.
type TypedOptionSetConfigProvider interface {
	ConfigProvider
	OptionTypeSet

	GetPlainOptionType() OptionType
	GetTypeOptionType() OptionType

	// IsExplicitlySelected returns, whether the provider is
	// selected by the given options.
	IsExplicitlySelected(opts Options) bool
}

// _TypedOptionSetConfigProvider is a private type used for private field embedding.
type _TypedOptionSetConfigProvider = TypedOptionSetConfigProvider

////////////////////////////////////////////////////////////////////////////////

type plainConfigProvider struct {
	ConfigOptionTypeSetHandler
}

var _ TypedOptionSetConfigProvider = (*plainConfigProvider)(nil)

// NewPlainConfigProvider create a TypedOptionSetConfigProvider without using
// a selective OptionType.
// The provided is selected if any option of the described OptionTypeSet has been
// given by the command line (no selective option type given).
// It can provide a Config from given Options.
func NewPlainConfigProvider(name string, adder ConfigAdder, types ...OptionType) TypedOptionSetConfigProvider {
	h := NewConfigOptionTypeSetHandler(name, adder, types...)
	return &plainConfigProvider{
		ConfigOptionTypeSetHandler: h,
	}
}

func (p *plainConfigProvider) GetConfigOptionTypeSet() OptionTypeSet {
	return p
}

func (p *plainConfigProvider) GetPlainOptionType() OptionType {
	return nil
}

func (p *plainConfigProvider) GetTypeOptionType() OptionType {
	return nil
}

func (p *plainConfigProvider) IsExplicitlySelected(opts Options) bool {
	return opts.FilterBy(p.HasOptionType).Changed()
}

func (p *plainConfigProvider) GetConfigFor(opts Options) (Config, error) {
	if !p.IsExplicitlySelected(opts) {
		return nil, nil
	}
	config := Config{}
	err := p.ApplyConfig(opts, config)
	return config, err
}

////////////////////////////////////////////////////////////////////////////////

type typedConfigProvider struct {
	_TypedOptionSetConfigProvider
	typeOptionType OptionType
}

var _ TypedOptionSetConfigProvider = (*typedConfigProvider)(nil)

// NewTypedConfigProvider creates a ConfigProvider using a selective
// type option (of type string) for choosing among multiple variants given as nested providers.
// The provider is selected if the option for this type or the plain option
// has been given on the command line.
func NewTypedConfigProvider(name string, desc, typeOption string, acceptUnknown ...bool) TypedOptionSetConfigProvider {
	if typeOption == "" {
		typeOption = name + "Type"
	}
	typeOpt := NewStringOptionType(typeOption, "type of "+desc)
	return &typedConfigProvider{NewTypedConfigProviderBase(name, desc, TypeNameProviderFromOptions(typeOption), general.Optional(acceptUnknown...), typeOpt), typeOpt}
}

func (p *typedConfigProvider) GetTypeOptionType() OptionType {
	return p.typeOptionType
}

func (p *typedConfigProvider) IsExplicitlySelected(opts Options) bool {
	return opts.Changed(p.typeOptionType.GetName(), p.GetPlainOptionType().GetName())
}

///////////////////////////////////////////////////////////////////////////////

// TypeNameProviderFromOptions offers a function extractiong
// a type name from the given option name. The Option's value
// must be a string. which is used a type name.
func TypeNameProviderFromOptions(name string) TypeNameProvider {
	return func(opts Options) (string, error) {
		typv, _ := opts.GetValue(name)
		typ, ok := typv.(string)
		if !ok {
			return "", fmt.Errorf("failed to assert type %T as string", typv)
		}
		return typ, nil
	}
}

type ExplicitlyTypedConfigTypeOptionSetConfigProvider interface {
	TypedOptionSetConfigProvider
	SetTypeName(n string)
}

type explicitlyTypedConfigProvider struct {
	_TypedOptionSetConfigProvider
	typeName string
}

var _ TypedOptionSetConfigProvider = (*typedConfigProvider)(nil)

// NewExplicitlyTypedConfigProvider provides a ConfigProvider
// using a fixed type name. OptionTypes and fixed type must be added separately.
func NewExplicitlyTypedConfigProvider(name string, desc string, acceptUnknown ...bool) ExplicitlyTypedConfigTypeOptionSetConfigProvider {
	p := &explicitlyTypedConfigProvider{}
	p._TypedOptionSetConfigProvider = NewTypedConfigProviderBase(name, desc, p.getTypeName, general.Optional(acceptUnknown...))
	return p
}

func (p *explicitlyTypedConfigProvider) SetTypeName(n string) {
	p.typeName = n
}

func (p *explicitlyTypedConfigProvider) getTypeName(opts Options) (string, error) {
	return p.typeName, nil
}

////////////////////////////////////////////////////////////////////////////////

type TypeNameProvider func(opts Options) (string, error)

type typedConfigProviderBase struct {
	OptionTypeSet
	typeProvider    TypeNameProvider
	meta            OptionTypeSet
	acceptUnknown   bool
	plainOptionType OptionType
	typefield       string
}

var _ TypedOptionSetConfigProvider = (*typedConfigProviderBase)(nil)

// NewTypedConfigProviderBase provides a base implementation for a ConfigProvider
// distinguishing among multiple config variants. Variants are given by nested
// OptionSets. Those sets must implement ConfigProvider to be able to
// // extract appropriate config. The actual variant is selected
// by a TypeNameProvider, typically from the actually given option settings.
// It uses an additional (plain) config option with the name of the provider
// accepting a structured value using
// a YAML option type. If this option is give , its value is used as Config value set.
// The type name might then be specified via the attribute `type`.
// It is selected if the TypeNameProvider is able to deliver a type name.
// If any Option of the set is given, the type setting is required, also.
func NewTypedConfigProviderBase(name string, desc string, prov TypeNameProvider, acceptUnknown bool, types ...OptionType) *typedConfigProviderBase {
	plainType := NewValueMapYAMLOptionType(name, desc+" (YAML)")
	set := NewOptionTypeSet(name, append(types, plainType)...)
	return &typedConfigProviderBase{
		OptionTypeSet:   set,
		typeProvider:    prov,
		meta:            NewOptionTypeSet(name, append(types, NewValueMapYAMLOptionType(name, desc+" (YAML)"))...),
		acceptUnknown:   acceptUnknown,
		plainOptionType: plainType,
		typefield:       "type",
	}
}

func (p *typedConfigProviderBase) WithTypeField(name string) *typedConfigProviderBase {
	p.typefield = name
	return p
}

func (p *typedConfigProviderBase) GetConfigOptionTypeSet() OptionTypeSet {
	return p
}

func (p *typedConfigProviderBase) GetPlainOptionType() OptionType {
	return p.plainOptionType
}

func (p *typedConfigProviderBase) GetTypeOptionType() OptionType {
	return nil
}

func (p *typedConfigProviderBase) IsExplicitlySelected(opts Options) bool {
	t, err := p.typeProvider(opts)
	return err == nil && t != ""
}

func (p *typedConfigProviderBase) GetConfigFor(opts Options) (Config, error) {
	typ, err := p.typeProvider(opts)
	if err != nil {
		return nil, err
	}
	cfgv, _ := opts.GetValue(p.GetName())

	var data Config
	if cfgv != nil {
		var ok bool
		data, ok = cfgv.(Config)
		if !ok {
			return nil, fmt.Errorf("failed to assert type %T as Config", cfgv)
		}
	}

	opts = opts.FilterBy(p.HasOptionType)
	if typ == "" && data != nil && data[p.typefield] != nil {
		t := data[p.typefield]
		if t != nil {
			if s, ok := t.(string); ok {
				typ = s
			} else {
				return nil, fmt.Errorf("type field must be a string")
			}
		}
	}

	if opts.Changed() || typ != "" {
		if typ == "" {
			return nil, fmt.Errorf("type required for explicitly configured options")
		}
		if data == nil {
			data = Config{}
		}
		data["type"] = typ
		if err := p.applyConfigForType(typ, opts, data); err != nil {
			if !p.acceptUnknown || !errors.Is(err, errors.ErrUnknown(typ)) {
				return nil, err
			}
			unexpected := opts.FilterBy(And(Changed(opts), Not(p.meta.HasOptionType))).Names()
			if len(unexpected) > 0 {
				return nil, fmt.Errorf("unexpected options %s", strings.Join(unexpected, ", "))
			}
		}
	}
	return data, nil
}

func (p *typedConfigProviderBase) GetTypeSetForType(name string) OptionTypeSet {
	set := p.GetTypeSet(name)
	if set == nil {
		k, v := KindVersion(name)
		if v == "" {
			set = p.GetTypeSet(TypeName(name, "v1"))
		} else if v == "v1" {
			set = p.GetTypeSet(k)
		}
	}
	return set
}

func (p *typedConfigProviderBase) applyConfigForType(name string, opts Options, config Config) error {
	set := p.GetTypeSetForType(name)
	if set == nil {
		return errors.ErrUnknown(name)
	}

	opts = opts.FilterBy(Not(p.meta.HasOptionType))
	err := opts.Check(set, p.GetName()+" type "+name)
	if err != nil {
		return err
	}
	handler, ok := set.(ConfigHandler)
	if !ok {
		return fmt.Errorf("no config handler defined for %s type %s", p.GetName(), name)
	}
	return handler.ApplyConfig(opts, config)
}

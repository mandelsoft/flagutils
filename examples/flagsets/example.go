package flagsets

import (
	"github.com/mandelsoft/flagutils/flagsets"
	"github.com/mandelsoft/flagutils/flagsets/scheme"
)

// This examples handles the creation of variations of some kind
// of typed object (here the base type Object) using commandline
// arguments to select the variant type and the attribution of the
// chosen variant.

// Object is the base type for all supported variants.
type Object interface {
	GetType() string
}

// Scheme describes the set of all possible variants.
var Scheme scheme.Scheme[Object]

// init configures the scheme with two implemented variants.
func init() {
	Scheme = scheme.New[Object]()
	Scheme.AddType(NewTypeA())
	Scheme.AddType(NewTypeB())
}

type ObjectMeta struct {
	Type string `json:"type"`
}

func (o *ObjectMeta) GetType() string {
	return o.Type
}

////////////////////////////////////////////////////////////////////////////////

// CommonOption describes an option type for an option for some configuratipon attribute
// which will be used by multiple variants.
var CommonOption = flagsets.NewStringOptionType("common", "Common Attribute")

////////////////////////////////////////////////////////////////////////////////

const TYPE_A = "typeA"

// TypeA is out first implementation variant with the type TYPE_A.
// It uses a field Common and AttrA.
type TypeA struct {
	ObjectMeta
	Common string `json:"common"`
	AttrA  string `json:"attrA"`
}

// AttrAOption is the option type to configure a value for AttrA.
var AttrAOption = flagsets.NewStringOptionType("attra", "Attribute A")

// NewTypeA creates a type object for the implementation variant typeA.
// It configured the type name, the mapping of options to a Config object
// and the options required to configure such a variant.
// This information is used by scheme.NewType to create a flagsets.TypedOptionSetConfigProvider.
// which is then used to compose a flagsets.TypedOptionSetConfigProvider.
func NewTypeA() scheme.Type[Object] {
	return scheme.NewType[TypeA, Object](
		TYPE_A, ConfigureTypeA,
		CommonOption,
		AttrAOption,
	)
}

// ConfigureTypeA transfers Options for TypeA to a flagsets.Config Object,
// which can be used to unmarshal an object of type TypeA.
func ConfigureTypeA(opts flagsets.Options, config flagsets.Config) error {
	// transfer option values to a structured confiig object.
	flagsets.AddFieldByOptionP(opts, CommonOption, config, "common")
	flagsets.AddFieldByOptionP(opts, AttrAOption, config, "attrA")
	return nil
}

////////////////////////////////////////////////////////////////////////////////

const TYPE_B = "typeB"

// TypeB is the same for a second implementation variant.
// Is also uses a field of the Common kind and an AttrB.
type TypeB struct {
	ObjectMeta
	Common string `json:"common"`
	AttrB  string `json:"attrB"`
}

// AttrBOption is the option type for an AttrB.
var AttrBOption = flagsets.NewStringOptionType("attrb", "Attribute B")

// NewTypeB creates a type object for the implementation variant typeB.
// It configures the type name, the mapping of options to a Config object
// and the options required to configure such a variant.
func NewTypeB() scheme.Type[Object] {
	return scheme.NewType[TypeB, Object](
		TYPE_B, ConfigureTypeB,
		CommonOption,
		AttrBOption,
	)
}

// ConfigureTypeB transfers Options for TypeB to a flagsets.Config Object,
// which can be used to unmarshal an object of type TypeA.
func ConfigureTypeB(opts flagsets.Options, config flagsets.Config) error {
	flagsets.AddFieldByOptionP(opts, CommonOption, config, "common")
	flagsets.AddFieldByOptionP(opts, AttrBOption, config, "attrB")
	return nil
}

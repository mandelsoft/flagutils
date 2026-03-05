package flagsets

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mandelsoft/flagutils/flagsets"
	"github.com/mandelsoft/goutils/funcs"
	"github.com/spf13/pflag"
)

func EvaluateArguments(provider flagsets.TypedOptionSetConfigProvider, args ...string) (flagsets.Config, error) {
	flags := pflag.NewFlagSet("flags", pflag.ContinueOnError)

	// the provider is used to create an explicit option set, which is then
	// used to feed the pflag.FlagSet,
	opts := provider.CreateOptions()
	opts.AddFlags(flags)

	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}

	// After parsing the command line options,
	// the provider is finally used to create a flagsets.Config object
	// for the implementation variant.
	return provider.GetConfigFor(opts)
}

func Usage() error {
	// We want to use command line options to configure different implementation
	// variants for some common interface. The variants are registered in the
	// scheme Scheme. This scheme uses the flagsets.OptionSet toolset
	// to configure possible overlapping OptionType sets for properties
	// of the known variants. Those sets are then combined to a
	// flagsets.TypedOptionSetConfigProvider.
	// It is used to create final options for a pflag.FlagSet, which is
	// then used to evaluate the given options to select and configure
	// the selected implementation variant.

	// The scheme provides a flagset.TypedOptionSetConfigProvider.
	// It is able to detect the variant by a type flag and appropriate
	// flags for this type.
	provider, err := Scheme.CreateOptionSetConfigProvider()
	if err != nil {
		return err
	}

	// we simulate command line flags, here.
	// the objectType flag is used to select the variant.
	// Additionally, the options for the required variant attributes are used
	// as described by the Type object for the variant.

	// if we use a wrong combination the problem is detected and reported
	cli := []string{
		"--attra=valueA",
		"--attrb=any", // this option describes a field not used by variant typeA
		"--common=valueCommon",
		"--objectType=typeA",
	}

	cfg, err := EvaluateArguments(provider, cli...)
	if err == nil {
		return fmt.Errorf("expected configuration problem not found")
	}
	if err.Error() != "option \"attrb\" given, but not possible for object type typeA" {
		return fmt.Errorf("expected error did not occur but found %w", err)
	}

	// now we omit the problematic option
	cli = []string{
		"--attra=valueA",
		"--common=valueCommon",
		"--objectType=typeA",
	}
	cfg, err = EvaluateArguments(provider, cli...)
	if err != nil {
		return err
	}

	fmt.Printf("found config: %+v\n", cfg)

	// This Config is then used by the Scheme to create
	// and configure the implementation object.
	a, err := Scheme.CreateObject(cfg)
	if err != nil {
		return err
	}

	expected := &TypeA{
		ObjectMeta: ObjectMeta{"typeA"},
		Common:     "valueCommon",
		AttrA:      "valueA",
	}

	if !reflect.DeepEqual(a, expected) {
		return fmt.Errorf("expected:\n%#v\ngot:\n%#v", expected, a)
	}

	fmt.Printf("result: %s", funcs.First(json.Marshal(a)))
	return nil
}

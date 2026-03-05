package flagsets_test

import (
	"github.com/mandelsoft/flagutils/examples/flagsets"
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/pflag"
)

var _ = Describe("variation test", func() {
	scheme := flagsets.Scheme

	handler := Must(scheme.CreateOptionSetConfigProvider())

	var flags *pflag.FlagSet

	BeforeEach(func() {
		flags = pflag.NewFlagSet("flags", pflag.ContinueOnError)
	})

	It("handles config of typeA", func() {
		opts := handler.CreateOptions()
		opts.AddFlags(flags)

		cli := []string{
			"--attra=valueA",
			"--common=valueCommon",
			"--objectType=typeA",
		}

		MustBeSuccessful(flags.Parse(cli))

		cfg := Must(handler.GetConfigFor(opts))

		a := Must(scheme.CreateObject(cfg))

		Expect(a).To(DeepEqual(&flagsets.TypeA{
			ObjectMeta: flagsets.ObjectMeta{"typeA"},
			Common:     "valueCommon",
			AttrA:      "valueA",
		}))
	})

	It("detects misconfig", func() {
		opts := handler.CreateOptions()
		opts.AddFlags(flags)

		cli := []string{
			"--attrb=valueB",
			"--common=valueCommon",
			"--objectType=typeA",
		}

		MustBeSuccessful(flags.Parse(cli))

		ExpectError(handler.GetConfigFor(opts)).To(MatchError(`option "attrb" given, but not possible for object type typeA`))
	})
})

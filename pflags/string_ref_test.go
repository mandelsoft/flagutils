package pflags

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/pflag"
)

var _ = Describe("string ref flags", func() {
	var flags *pflag.FlagSet

	BeforeEach(func() {
		flags = pflag.NewFlagSet("test", pflag.ContinueOnError)
	})

	It("handles string", func() {
		var flag *string
		StringRefVarP(flags, &flag, "flag", "", nil, "test flag")

		value := `value`

		Expect(flags.Parse([]string{"--flag", value})).To(Succeed())
		Expect(flag).NotTo(BeNil())
		Expect(*flag).To(Equal("value"))

		Expect(GetStringRef(flags, "flag")).To(Equal(flag))
	})

	It("keeps nil if not set", func() {
		var flag *string
		StringRefVarP(flags, &flag, "flag", "", nil, "test flag")

		err := flags.Parse([]string{})
		Expect(err).NotTo(HaveOccurred())
		Expect(flag).To(BeNil())
	})
})

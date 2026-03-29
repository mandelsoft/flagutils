package pflags

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/pflag"
)

var _ = Describe("float ref flags", func() {
	var flags *pflag.FlagSet

	BeforeEach(func() {
		flags = pflag.NewFlagSet("test", pflag.ContinueOnError)
	})

	It("handles string", func() {
		var flag *float64
		Float64RefVarP(flags, &flag, "flag", "", nil, "test flag")

		value := `3.14`

		Expect(flags.Parse([]string{"--flag", value})).To(Succeed())
		Expect(flag).NotTo(BeNil())
		Expect(*flag).To(Equal(float64(3.14)))

		Expect(GetFloat64Ref(flags, "flag")).To(Equal(flag))
	})

	It("keeps nil if not set", func() {
		var flag *float64
		Float64RefVarP(flags, &flag, "flag", "", nil, "test flag")

		err := flags.Parse([]string{})
		Expect(err).NotTo(HaveOccurred())
		Expect(flag).To(BeNil())
	})
})

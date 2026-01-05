package flagutils_test

import (
	"context"
	"github.com/mandelsoft/flagutils"
	"github.com/spf13/pflag"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Nesting interface {
	Special()
}
type nesting struct {
	options flagutils.DefaultOptionSet
	testOpt TestOption
}

func NewNesting() *nesting {
	m := &nesting{}
	m.options.Add(&m.testOpt)
	return m
}

var (
	_ flagutils.Options           = (*nesting)(nil)
	_ flagutils.OptionSetProvider = (*nesting)(nil)
)

func (n *nesting) AsOptionSet() flagutils.OptionSet {
	return n.options.AsOptionSet()
}

func (n *nesting) AddFlags(fs *pflag.FlagSet) {
}

func (n *nesting) Special() {
}

var _ = Describe("options", func() {
	var set flagutils.DefaultOptionSet

	BeforeEach(func() {
		set = nil
	})

	Context("option set", func() {
		It("validates", func() {
			n := NewNesting()
			set.Add(n)
			Expect(flagutils.Validate(context.Background(), set, nil)).To(Succeed())
			Expect(n.testOpt.Validated).To(BeTrue())
		})

		It("retrieves nested", func() {
			set.Add(NewNesting())
			opt := flagutils.GetFrom[*TestOption](set)
			Expect(opt).NotTo(BeNil())
		})

		It("retrieves interface", func() {
			set.Add(NewNesting())
			opt := flagutils.GetFrom[Nesting](set)
			Expect(opt).NotTo(BeNil())
		})
	})
})

package flagutils_test

import (
	"context"
	"strings"

	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/goutils/iterutils"
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
	var set flagutils.ExtendableOptionSet

	BeforeEach(func() {
		set = flagutils.NewOptionSet()
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

	Context("assure", func() {
		It("adds", func() {
			flagutils.Assure(set, NewTestOption())
			Expect(len(iterutils.Get(set.Options))).To(Equal(1))
			Expect(flagutils.GetFrom[*TestOption](set)).NotTo(BeNil())
		})

		It("keeps", func() {
			set.Add(NewTestOption()())
			flagutils.Assure(set, NewTestOption())
			Expect(len(iterutils.Get(set.Options))).To(Equal(1))
			Expect(flagutils.GetFrom[*TestOption](set)).NotTo(BeNil())
		})

		It("adds filtered", func() {
			set.Add(NewTestOption("one")())
			flagutils.Assure(set, NewTestOption("two"), check("two"))
			Expect(len(iterutils.Get(set.Options))).To(Equal(2))
			Expect(flagutils.Filter[*TestOption](set, check("two"))).NotTo(BeNil())
			Expect(flagutils.Filter[*TestOption](set, check("one"))).NotTo(BeNil())
		})

		It("keeps filtered", func() {
			set.Add(NewTestOption("one")())
			flagutils.Assure(set, NewTestOption("one"), check("one"))
			Expect(len(iterutils.Get(set.Options))).To(Equal(1))
			Expect(flagutils.Filter[*TestOption](set, check("one"))).NotTo(BeNil())
		})
	})
})

func check(mode ...string) func(o *TestOption) bool {
	m := strings.Join(mode, "")
	return func(o *TestOption) bool {
		return o.Mode == m
	}
}

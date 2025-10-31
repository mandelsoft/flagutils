package flagutils_test

import (
	"github.com/mandelsoft/flagutils"
	"github.com/spf13/pflag"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Interface interface {
	Get() bool
}

type TestOption struct {
	Flag bool
}

func (t *TestOption) Get() bool {
	return t.Flag
}

func (t *TestOption) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&t.Flag, "test", "t", false, "test flag")
}

type Test2Option struct {
	Flag bool
}

func (t *Test2Option) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&t.Flag, "flag", "f", false, "flag")
}

type ParamOption[I any] struct {
	Flag bool
}

func (t *ParamOption[I]) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&t.Flag, "flag", "f", false, "flag")
}

type SetOption struct {
	flagutils.DefaultOptionSet
	Flag bool
}

var _ flagutils.Options = (*TestOption)(nil)

var _ = Describe("options", func() {
	var set flagutils.DefaultOptionSet

	BeforeEach(func() {
		set = nil
	})

	Context("simple option", func() {
		It("skips unknown option", func() {
			var opt *TestOption
			Expect(flagutils.RetrieveFrom(set, opt)).To(BeFalse())
		})

		It("assigns options pointer from set", func() {
			var t2 *Test2Option

			inst := &TestOption{}
			set.Add(inst)
			flagutils.GetFrom[*TestOption](set).Flag = true

			var opt *TestOption
			Expect(flagutils.RetrieveFrom(set, &opt)).To(BeTrue())
			Expect(opt.Flag).To(BeTrue())
			Expect(opt).To(BeIdenticalTo(inst))

			Expect(flagutils.RetrieveFrom(set, &t2)).To(BeFalse())
		})

		It("assigns options value from set", func() {
			inst := &TestOption{}
			set.Add(inst)
			flagutils.GetFrom[*TestOption](set).Flag = true

			var opt TestOption
			Expect(flagutils.RetrieveFrom(set, &opt)).To(BeTrue())
			Expect(opt.Flag).To(BeTrue())

			opt.Flag = false
			Expect(inst.Flag).To(BeTrue())
		})
	})

	Context("set option", func() {
		It("skips unknown option", func() {
			var opt *SetOption
			Expect(flagutils.RetrieveFrom(set, opt)).To(BeFalse())
		})

		It("assigns options pointer from set", func() {
			var t2 *Test2Option

			inst := &SetOption{}
			set.Add(inst)
			flagutils.GetFrom[*SetOption](set).Flag = true

			var opt *SetOption
			Expect(flagutils.RetrieveFrom(set, &opt)).To(BeTrue())
			Expect(opt.Flag).To(BeTrue())
			Expect(opt).To(BeIdenticalTo(inst))

			Expect(flagutils.RetrieveFrom(set, &t2)).To(BeFalse())
		})

		It("assigns options value from set", func() {
			inst := &SetOption{}
			set.Add(inst)
			flagutils.GetFrom[*SetOption](set).Flag = true

			var opt SetOption
			Expect(flagutils.RetrieveFrom(set, &opt)).To(BeTrue())
			Expect(opt.Flag).To(BeTrue())

			opt.Flag = false
			Expect(inst.Flag).To(BeTrue())
		})
	})

	Context("nested option", func() {
		var group *SetOption

		BeforeEach(func() {
			group = &SetOption{}
			set.Add(group)
		})

		It("skips unknown option", func() {
			var opt *TestOption
			Expect(flagutils.RetrieveFrom(set, opt)).To(BeFalse())
		})

		It("assigns options pointer from set", func() {
			var t2 *Test2Option

			inst := &TestOption{}
			group.Add(inst)
			flagutils.GetFrom[*TestOption](set).Flag = true

			var opt *TestOption
			Expect(flagutils.RetrieveFrom(set, &opt)).To(BeTrue())
			Expect(opt.Flag).To(BeTrue())
			Expect(opt).To(BeIdenticalTo(inst))

			Expect(flagutils.RetrieveFrom(set, &t2)).To(BeFalse())
		})

		It("assigns options value from set", func() {
			inst := &TestOption{}
			group.Add(inst)
			flagutils.GetFrom[*TestOption](set).Flag = true

			var opt TestOption
			Expect(flagutils.RetrieveFrom(set, &opt)).To(BeTrue())
			Expect(opt.Flag).To(BeTrue())

			opt.Flag = false
			Expect(inst.Flag).To(BeTrue())
		})
	})

	Context("interface", func() {
		It("assigns options from set", func() {
			inst := &TestOption{}
			set.Add(inst)
			flagutils.GetFrom[*TestOption](set).Flag = true

			var opt Interface
			Expect(flagutils.RetrieveFrom(set, &opt)).To(BeTrue())
			Expect(opt.Get()).To(BeTrue())

			inst.Flag = false
			Expect(opt.Get()).To(BeFalse())
		})
	})

	Context("param option", func() {
		It("gets parameterized option", func() {
			inst := &ParamOption[int]{}
			set.Add(inst)

			var opt *ParamOption[int]
			Expect(flagutils.RetrieveFrom(set, &opt)).To(BeTrue())
		})
	})
})

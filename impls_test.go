package flagutils_test

import (
	"context"
	"github.com/mandelsoft/flagutils"
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
)

var _ = Describe("traversing", func() {

	Context("simple option", func() {
		It("default option set", func() {
			set := flagutils.DefaultOptionSet{}
			set.Add(flagutils.NoOptions{})
			set.Add(&TestOption{})

			MustBeSuccessful(flagutils.Validate(context.Background(), set, nil))
			o := flagutils.GetFrom[*TestOption](set)
			Expect(o).NotTo(BeNil())
			Expect(o.Validated).To(BeTrue())

			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			set.AddFlags(fs)
			MustBeSuccessful(fs.Parse([]string{"-t"}))
			Expect(o.Flag).To(BeTrue())
			Expect(o.Finalized).To(BeFalse())

			MustBeSuccessful(flagutils.Finalize(context.Background(), set, nil))
			Expect(o).NotTo(BeNil())
			Expect(o.Finalized).To(BeTrue())
		})

		It("option set", func() {
			set := &flagutils.SetBasedOptions{}
			set.Add(flagutils.NoOptions{})
			set.Add(&TestOption{})

			MustBeSuccessful(flagutils.Validate(context.Background(), set, nil))
			o := flagutils.GetFrom[*TestOption](set)
			Expect(o).NotTo(BeNil())
			Expect(o.Validated).To(BeTrue())
			Expect(o.Finalized).To(BeFalse())

			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			set.AddFlags(fs)
			MustBeSuccessful(fs.Parse([]string{"-t"}))
			Expect(o.Flag).To(BeTrue())

			MustBeSuccessful(flagutils.Finalize(context.Background(), set, nil))
			Expect(o).NotTo(BeNil())
			Expect(o.Finalized).To(BeTrue())
		})
	})
})

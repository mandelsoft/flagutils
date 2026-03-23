package flagutils_test

import (
	"context"
	"fmt"

	"github.com/mandelsoft/flagutils"
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
)

type Counter = int
type Countable interface {
	GetCount() Counter
}

func incCounter(ctx context.Context) int {
	c := ctx.Value("counter").(*Counter)
	*c++
	return *c
}

func withCounter(ctx context.Context) context.Context {
	var c Counter
	return context.WithValue(ctx, "counter", &c)
}

func Check[T Countable](set flagutils.OptionSet, n int) {
	o := flagutils.GetFrom[T](set)
	ExpectWithOffset(1, o).ToNot(BeNil())
	ExpectWithOffset(1, o.GetCount()).To(Equal(n))
}

type O1 struct {
	handled int
}

var _ flagutils.Preparable = (*O1)(nil)

var _ flagutils.Validatable = (*O1)(nil)
var _ flagutils.Options = (*O1)(nil)

func (o *O1) GetCount() Counter {
	return o.handled
}

func (o *O1) AddFlags(fs *pflag.FlagSet) {
}

func (o *O1) Prepare(ctx context.Context, opts flagutils.OptionSet, v flagutils.PreparationSet) error {
	if o.handled != 0 {
		return fmt.Errorf("already handled")
	}
	o.handled = incCounter(ctx)
	return nil
}

func (o *O1) Validate(ctx context.Context, opts flagutils.OptionSet, v flagutils.ValidationSet) error {
	if o.handled != 0 {
		return fmt.Errorf("already handled")
	}
	o.handled = incCounter(ctx)
	return nil
}

func (o *O1) Finalize(ctx context.Context, opts flagutils.OptionSet, v flagutils.FinalizationSet) error {
	if o.handled != 0 {
		return fmt.Errorf("already handled")
	}
	o.handled = incCounter(ctx)
	return nil
}

type O2 struct {
	O1
}

func (o *O2) Prepare(ctx context.Context, opts flagutils.OptionSet, v flagutils.PreparationSet) error {
	o1, err := flagutils.PreparedOptions[*O1](ctx, opts, v)
	if err != nil {
		return err
	}
	if o1 == nil {
		return fmt.Errorf("o1 not found")
	}
	return o.O1.Prepare(ctx, opts, v)
}

func (o *O2) Finalize(ctx context.Context, opts flagutils.OptionSet, v flagutils.FinalizationSet) error {
	err := o.O1.Finalize(ctx, opts, v)
	if err != nil {
		return err
	}
	o1, err := flagutils.FinalizedOptions[*O1](ctx, opts, v)
	if o1 == nil {
		return fmt.Errorf("o1 not found")
	}
	return err
}

type O3 struct {
	O1
}

type O4 struct {
	O1
	flagutils.DefaultOptionSet
}

func newO4() *O4 {
	o := &O4{}
	o.Add(&O3{})
	return o
}

func (o *O4) AddFlags(flag *pflag.FlagSet) {
}

func (o *O4) AsOptionSet() flagutils.OptionSet {
	return o.DefaultOptionSet
}

func (o *O4) Prepare(ctx context.Context, opts flagutils.OptionSet, v flagutils.PreparationSet) error {
	err := v.PrepareSet(ctx, opts, o.AsOptionSet())
	if err != nil {
		return err
	}
	return o.O1.Prepare(ctx, opts, v)
}

func (o *O4) Finalize(ctx context.Context, opts flagutils.OptionSet, v flagutils.FinalizationSet) error {
	err := o.O1.Finalize(ctx, opts, v)
	if err != nil {
		return err
	}
	return v.FinalizeSet(ctx, opts, o.AsOptionSet())
}

var _ = Describe("Lifecycle Test Environment", func() {
	var set flagutils.DefaultOptionSet
	var ctx context.Context

	BeforeEach(func() {
		ctx = withCounter(context.Background())
		set = flagutils.DefaultOptionSet{}
		set.Add(&O2{})
		set.Add(newO4())
		set.Add(&O1{})
	})

	Context("Prepare", func() {
		It("ordered", func() {
			MustBeSuccessful(flagutils.Prepare(ctx, set, nil))

			Check[*O1](set, 1)
			Check[*O2](set, 2)
			Check[*O3](set, 3)
			Check[*O4](set, 4)
		})
	})

	Context("Finalizes", func() {
		It("ordered", func() {
			MustBeSuccessful(flagutils.Finalize(ctx, set, nil))

			Check[*O1](set, 1)
			Check[*O2](set, 4)
			Check[*O4](set, 2)
			Check[*O3](set, 3)
		})
	})
})

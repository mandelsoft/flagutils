package flagsets_test

import (
	"github.com/mandelsoft/flagutils/examples/flagsets"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Example Test Environment", func() {
	It("usage", func() {
		Expect(flagsets.Usage()).To(Succeed())
	})
})

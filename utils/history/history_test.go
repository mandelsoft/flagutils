package history_test

import (
	"github.com/mandelsoft/flagutils/utils/history"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"slices"
	"strings"
)

func H(s string) history.History[string] {
	return history.History[string](strings.Split(s, "/"))
}

var _ = Describe("History", func() {
	cmp := history.CompareFunc[string](strings.Compare)

	Context("compare", func() {
		It("should compare paths correctly", func() {
			Expect(cmp(H("a/b"), H("a/b"))).To(Equal(0))
			Expect(cmp(H("a"), H("a/b"))).To(Equal(-1))
			Expect(cmp(H("a/b"), H("b"))).To(Equal(-1))

			Expect(cmp(H("b"), H("a/b"))).To(Equal(1))
		})
	})

	Context("sort", func() {
		It("sorts correctly", func() {
			hist := []history.History[string]{H("a"), H("b"), H("c"), H("c/a"), H("c/b"), H("b/b"), H("b/a")}
			slices.SortStableFunc(hist, cmp)
			Expect(hist).To(Equal([]history.History[string]{H("a"), H("b"), H("b/a"), H("b/b"), H("c"), H("c/a"), H("c/b")}))
		})
	})
})

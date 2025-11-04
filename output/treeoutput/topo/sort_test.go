package topo_test

import (
	"fmt"
	"github.com/mandelsoft/goutils/general"
	"slices"
	"strings"

	"github.com/mandelsoft/flagutils/output/treeoutput/topo"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Element struct {
	topo.TopoInfo[string, string]

	index int
}

func (e *Element) String() string {
	return fmt.Sprintf("%s(%d)", strings.Join(e.GetHierarchy(), "/"), e.index)
}

func IdProvider(path []string) string {
	return topo.StringIdProviderFunc(path)
}

func NewElement(key string, index int, path ...string) *Element {
	return &Element{index: index, TopoInfo: topo.NewStringIdTopoInfo[string](key, path)}
}

var _ = Describe("TopoSorter", func() {
	var elements []*Element
	var cmp general.CompareFunc[*Element]

	A := NewElement("a", 14)
	AD := NewElement("d", 13, "a")
	ADJ := NewElement("j", 11, "a", "d")
	ADK := NewElement("k", 12, "a", "d")
	AE := NewElement("e", 10, "a")
	B := NewElement("b", 9)
	BF := NewElement("f", 8, "b")
	BFL := NewElement("l", 7, "b", "f")
	BFM := NewElement("m", 6, "b", "f")
	BG := NewElement("g", 5, "b")
	C := NewElement("c", 4)
	CH := NewElement("h", 3, "c")
	CI := NewElement("i", 2, "c")
	CIN := NewElement("n", 1, "c", "i")
	CIO := NewElement("o", 0, "c", "i")

	BeforeEach(func() {
		elements = []*Element{
			/*  0: c/i/o */ CIO,
			/*  1: c/i/n */ CIN,
			/*  2; c/i   */ CI,
			/*  3: c/h   */ CH,
			/*  4: c     */ C,
			/*  5: b/g   */ BG,
			/*  6: b/f/m */ BFM,
			/*  7: b/f/l */ BFL,
			/*  8: b/f   */ BF,
			/*  9: b     */ B,
			/* 10: a/e   */ AE,
			/* 11: a/d/j */ ADJ,
			/* 12: a/d/k */ ADK,
			/* 13: a/d   */ AD,
			/* 14: a     */ A,
		}
		cmp = topo.NewDefaultCompareFunc(elements, IdProvider)

	})

	Context("Compare", func() {
		It("compare", func() {
			Expect(cmp(A, B)).To(BeNumerically(">", 0))
			Expect(cmp(B, C)).To(BeNumerically(">", 0))
			Expect(cmp(BF, C)).To(BeNumerically(">", 0))
			Expect(cmp(BFL, C)).To(BeNumerically(">", 0))
			Expect(cmp(BF, BG)).To(BeNumerically(">", 0))
			Expect(cmp(BFL, BG)).To(BeNumerically(">", 0))
			Expect(cmp(ADJ, BFM)).To(BeNumerically(">", 0))
			Expect(cmp(ADJ, AD)).To(BeNumerically(">", 0))
			Expect(cmp(ADJ, ADK)).To(BeNumerically("<", 0))
		})
	})

	Context("Sort", func() {
		It("Sort", func() {
			slices.SortFunc(elements, cmp)
			Expect(elements).To(Equal([]*Element{C, CI, CIO, CIN, CH, B, BG, BF, BFM, BFL, A, AE, AD, ADJ, ADK}))
		})
	})
})

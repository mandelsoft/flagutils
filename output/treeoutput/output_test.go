package treeoutput_test

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/closure"
	"github.com/mandelsoft/flagutils/examples/files/files"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/flagutils/output/tableoutput"
	"github.com/mandelsoft/flagutils/sort"
	"github.com/mandelsoft/flagutils/utils/out"

	"github.com/spf13/pflag"

	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree Output", func() {
	var ctx context.Context
	var opts flagutils.DefaultOptionSet
	var fs *pflag.FlagSet
	var outp *bytes.Buffer

	BeforeEach(func() {
		outp = bytes.NewBuffer(nil)
		ctx = out.With(context.Background(), out.New(outp, os.Stderr))
		opts = flagutils.DefaultOptionSet{}
		opts.Add(
			files.New(),
			closure.NewByFactory[*files.Element](files.ClosureFactory),
			sort.New(),
			tableoutput.New(),
			output.New(files.OutputsFactory),
		)
		fs = pflag.NewFlagSet("test", pflag.ContinueOnError)
		opts.AddFlags(fs)
	})

	AfterEach(func() {
		flagutils.Finalize(ctx, opts, nil)
	})

	Context("when creating tree output", func() {
		It("should initialize correctly", func() {
			MustBeSuccessful(fs.Parse([]string{"-c", "-s", "name", "output", "test", "-o", "tree"}))
			MustBeSuccessful(flagutils.Validate(ctx, opts, nil))
			args := fs.Args()
			out := output.From[*files.Element](opts)
			n := Must(out.GetOutput().Process(ctx, args, files.NewSourceFactory(opts)))
			Expect(n).To(Equal(10))

			exp := `
            MODE       NAME   SIZE ERROR
├─                     output      .*file.*
└─ ⊗        drwxrwxr.x test   *\d+ 
   ├─       -rw-rw-r.- a         5 
   ├─       -rw-rw-r.- b         3 
   └─ ⊗     drwxrwxr.x dir    *\d+ 
      ├─    -rw-rw-r.- a         5 
      ├─    -rw-rw-r.- c         6 
      └─ ⊗  drwxrwxr.x sub    *\d+ 
         ├─ -rw-rw-r.- d         6 
         └─ -rw-rw-r.- e         3 
`
			Expect(outp.String()).To(StringMatchTrimmedWithContext(exp))
		})
	})
})

func compareRunes(a, b string) string {
	line := 1
	s := ""
	ra, rb := []rune(a), []rune(b)
	for i := range ra {
		if ra[i] == '\n' {
			line++
			s = ""
		} else {
			s += string(ra[i])
		}
		if len(rb) < i {
			return fmt.Sprintf("additional rune %c", ra)
		} else {
			if ra[i] != rb[i] {
				return fmt.Sprintf("different rune %c (expected %c) (line %d: %s)", ra[i], rb[i], line, s)
			}
		}
	}
	if len(ra) < len(rb) {
		return fmt.Sprintf("missing rune %c", rb[len(ra)])
	}
	return ""
}

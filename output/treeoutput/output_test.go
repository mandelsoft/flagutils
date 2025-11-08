package treeoutput_test

import (
	"bytes"
	"context"
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
			Expect(outp.String()).To(StringEqualTrimmedWithContext(`
            MODE       NAME   SIZE ERROR
├─                     output      GetFileAttributesEx output: The system cannot find the file specified.
└─ ⊗        drwxrwxrwx test      0 
   ├─       -rw-rw-rw- a         5 
   ├─       -rw-rw-rw- b         3 
   └─ ⊗     drwxrwxrwx dir       0 
      ├─    -rw-rw-rw- a         5 
      ├─    -rw-rw-rw- c         6 
      └─ ⊗  drwxrwxrwx sub       0 
         ├─ -rw-rw-rw- d         6 
         └─ -rw-rw-rw- e         3 
`))
		})
	})
})

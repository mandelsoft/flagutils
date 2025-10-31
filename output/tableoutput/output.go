package tableoutput

import (
	"context"
	"fmt"
	"github.com/mandelsoft/flagutils/out"
	"github.com/mandelsoft/flagutils/output"
	"github.com/mandelsoft/goutils/iterutils"
	"github.com/mandelsoft/streaming"
	"iter"
	"strings"
)

type Factory struct {
	Headers []string
	Options *Options
}

var _ streaming.ProcessorFactory[output.ElementSpecs, int, []string] = (*Factory)(nil)

func (o *Factory) Processor(output.ElementSpecs) (streaming.Processor[int, []string], error) {
	return newProcessor(o).Process, nil
}

type Processor struct {
	output *Factory
	data   [][]string
}

var (
	_ streaming.Processor[int, []string] = (*Processor)(nil).Process
)

func newProcessor(o *Factory) *Processor {
	return &Processor{
		output: o,
	}

}

func (p *Processor) Process(ctx context.Context, i iter.Seq[[]string]) (int, error) {
	p.data = iterutils.Get(i)

	if len(p.data) == 0 {
		out.Print(ctx, "no elements found\n")
		return 0, nil
	}
	effheader := p.output.Headers
	if p.output.Options.UseColumnOptimization() {
		effheader = p.optimizeColumns()
	}
	FormatTable(ctx, "", append([][]string{effheader}, p.data...))
	return len(p.data), nil
}

func (p *Processor) optimizeColumns() []string {
	headers := p.output.Headers
	if len(p.data) < 2 {
		return headers
	}
	cnt := p.output.Options.GetOptimizedColumns()

columns:
	for cnt > 0 && len(headers) > 1 {
		e := p.data[0]
		if len(e) <= 1 {
			break
		}
		v := e[0]
		for j := range p.data {
			e = p.data[j]
			if len(e) < 1 || e[0] != v {
				break columns
			}
		}
		// all row value identical, skip column
		headers = headers[1:]
		for j := range p.data {
			p.data[j] = p.data[j][1:]
		}
		cnt--
	}
	return headers
}

func FormatTable(ctx context.Context, gap string, data [][]string) {
	columns := []int{}
	maxLen := 0
	maxTitle := 0

	formats := []string{}
	if len(data) > 1 {
		for i, f := range data[0] {
			if strings.HasPrefix(f, "-") {
				formats = append(formats, "")
				data[0][i] = f[1:]
			} else {
				formats = append(formats, "-")
			}
			if len(data[0][i]) > maxTitle {
				maxTitle = len(data[0][i])
			}
		}
	}

	for _, row := range data {
		for i, col := range row {
			l := len([]rune(col))
			if i >= len(columns) {
				columns = append(columns, l)
			} else if columns[i] < l {
				columns[i] = l
			}
			if l > maxLen {
				maxLen = l
			}
		}
	}

	if len(columns) > 2 && maxLen > 200 {
		first := []string{}
		setSep := false
		for i, row := range data {
			if i == 0 {
				first = row
			} else {
				for c, col := range row {
					if c < len(first) {
						out.Printf(ctx, "%s%-*s: %s\n", gap, maxTitle, first[c], col)
					} else {
						out.Printf(ctx, "%s%d: %s\n", gap, c, col)
					}
					setSep = true
				}
				if setSep {
					out.Printf(ctx, "---\n")
					setSep = false
				}
			}
		}
	} else {
		format := gap
		for i, col := range columns {
			f := "-"
			if i < len(formats) {
				f = formats[i]
			}
			if i == len(columns)-1 && f == "-" {
				format = fmt.Sprintf("%s%%s ", format)
			} else {
				format = fmt.Sprintf("%s%%%s%ds ", format, f, col)
			}
		}
		format = format[:len(format)-1] + "\n"
		for _, row := range data {
			if len(row) > 0 {
				r := []interface{}{}
				for i := 0; i < len(columns); i++ {
					if i < len(row) {
						r = append(r, row[i])
					} else {
						r = append(r, "")
					}
				}
				out.Printf(ctx, format, r...)
			}
		}
	}
}

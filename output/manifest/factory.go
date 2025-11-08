package manifest

import (
	"context"
	"github.com/mandelsoft/flagutils/closure"
	"github.com/mandelsoft/streaming/chain"

	"github.com/mandelsoft/flagutils"
	output "github.com/mandelsoft/flagutils/output/internal"
)

type OutputFactory[I any] struct {
	formatter Formatter
}

var _ output.OutputFactory[int] = (*OutputFactory[int])(nil)

func NewOutputFactory[I any](formatter Formatter) *OutputFactory[I] {
	return &OutputFactory[I]{formatter}
}

func (o *OutputFactory[I]) GetFieldNames(string) []string {
	return nil
}

func (o *OutputFactory[I]) Create(ctx context.Context, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (output.Output[I], error) {
	c := chain.New[I]()
	e := closure.From[I](opts)
	if e != nil {
		f := e.GetExploderFactory(opts)
		if f != nil {
			c = chain.AddExplodeByFactory[I](c, f)
		}
	}
	return output.NewOutput[I, I](c, &Factory[I]{o}), nil
}

func AddManifestOutputs[I any](out output.OutputsFactory[I]) output.OutputsFactory[I] {
	out.Add("yaml", NewYAMLFactory[I](false))
	out.Add("YAML", NewYAMLFactory[I](true))
	out.Add("json", NewJSONFactory[I](false))
	out.Add("JSON", NewJSONFactory[I](true))
	return out
}

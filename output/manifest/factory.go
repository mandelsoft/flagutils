package manifest

import (
	"context"
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

func (o *OutputFactory[I]) GetFieldNames() []string {
	return nil
}

func (o *OutputFactory[I]) Create(ctx context.Context, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (output.Output[I], error) {
	return output.NewOutput[I, I](chain.New[I](), &Factory[I]{o}), nil
}

func AddManifestOutputs[I any](out output.OutputsFactory[I]) output.OutputsFactory[I] {
	out.Add("yaml", NewYAMLFactory[I](false))
	out.Add("YAML", NewYAMLFactory[I](true))
	out.Add("json", NewJSONFactory[I](false))
	out.Add("JSON", NewJSONFactory[I](true))
	return out
}

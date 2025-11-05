package output

import (
	"context"
	"fmt"
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/flagutils/output/manifest"
	"github.com/mandelsoft/goutils/maputils"
)

////////////////////////////////////////////////////////////////////////////////

const FIELD_MODE_OUTPUT = "<output>"

type outputsFactory[I any] struct {
	modes map[string]OutputFactory[I]
}

func NewOutputsFactory[I any](alt ...map[string]OutputFactory[I]) OutputsFactory[I] {
	m := make(map[string]OutputFactory[I])
	for _, i := range alt {
		for k, v := range i {
			m[k] = v
		}
	}
	return &outputsFactory[I]{modes: m}
}

func (f *outputsFactory[I]) GetModes() []string {
	return maputils.OrderedKeys(f.modes)
}

func (f *outputsFactory[I]) Add(mode string, out OutputFactory[I]) OutputsFactory[I] {
	f.modes[mode] = out
	return f
}

func (f *outputsFactory[I]) AddManifestOutputs() OutputsFactory[I] {
	return manifest.AddManifestOutputs(f)
}

func (f *outputsFactory[I]) GetFieldNames(mode, stage string) []string {
	of := f.modes[mode]
	if of == nil {
		return nil
	}
	return of.GetFieldNames(stage)
}

func (f *outputsFactory[I]) CreateOutput(ctx context.Context, mode string, opts flagutils.OptionSetProvider, v flagutils.ValidationSet) (Output[I], error) {
	of := f.modes[mode]
	if of == nil {
		return nil, fmt.Errorf("invalid output mode: %s", mode)
	}
	return of.Create(ctx, opts, v)
}

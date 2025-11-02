package manifest

import (
	"bytes"
	"context"
	"encoding/json"
	output "github.com/mandelsoft/flagutils/output/internal"
	"github.com/mandelsoft/flagutils/utils/out"
	"go.yaml.in/yaml/v3"
)

type Formatter interface {
	Format(ctx context.Context, values []Manifest) error
}

type Manifest interface {
	AsManifest() interface{}
}

////////////////////////////////////////////////////////////////////////////////

type ItemList struct {
	Items []interface{} `json:"items"`
}

func format(ctx context.Context, values []Manifest, formatter func(data any) ([]byte, error)) ([]byte, error) {
	items := &ItemList{}
	for _, m := range values {
		items.Items = append(items.Items, m.AsManifest())
	}
	return formatter(items)
}

////////////////////////////////////////////////////////////////////////////////

type YAML struct {
	docs bool
}

var _ Formatter = (*YAML)(nil)

func NewYAML(docs bool) *YAML {
	return &YAML{docs}
}

func NewYAMLFactory[I any](docs bool) output.OutputFactory[I] {
	return NewOutputFactory[I](NewYAML(docs))
}

func (f *YAML) Format(ctx context.Context, values []Manifest) error {
	if f.docs {
		d, err := format(ctx, values, yaml.Marshal)
		if err != nil {
			return err
		}
		_, err = out.Write(ctx, d)
		return err
	} else {
		for _, m := range values {
			_, err := out.Printf(ctx, "---\n")
			if err != nil {
				return err
			}
			d, err := yaml.Marshal(m.AsManifest())
			if err != nil {
				return err
			}
			_, err = out.Write(ctx, d)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type JSON struct {
	pretty bool
}

var _ Formatter = (*JSON)(nil)

func NewJSON(pretty bool) *JSON {
	return &JSON{pretty}
}

func NewJSONFactory[I any](pretty bool) output.OutputFactory[I] {
	return NewOutputFactory[I](NewJSON(pretty))
}

func (f *JSON) Format(ctx context.Context, values []Manifest) error {
	d, err := format(ctx, values, json.Marshal)
	if err != nil {
		return err
	}
	if f.pretty {
		var buf bytes.Buffer
		err = json.Indent(&buf, d, "", "  ")
		if err != nil {
			return err
		}
		err = buf.WriteByte('\n')
		if err != nil {
			return err
		}
		d = buf.Bytes()
	}
	_, err = out.Write(ctx, d)
	return err
}

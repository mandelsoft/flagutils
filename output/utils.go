package output

import (
	"context"
	"fmt"
	"github.com/mandelsoft/flagutils"
	"github.com/mandelsoft/goutils/sliceutils"
	"sort"
)

func CheckFieldNames(fieldnames []string, ctx context.Context, v flagutils.ValidationSet, opts flagutils.OptionSet) error {
	fields, err := flagutils.ValidatedOptions[FieldNameProvider](ctx, opts, v)
	if err != nil {
		return err
	}

	if fields == nil {
		return fmt.Errorf("invalid fields: %v", fieldnames)
	}
	names := fields.GetFieldNames()
	if names == nil {
		return fmt.Errorf("invalid fields: %v", fieldnames)
	}

	wrong := sliceutils.Diff(fieldnames, names)
	if len(wrong) != 0 {
		sort.Strings(wrong)
		return fmt.Errorf("invalid sort fields: %v", wrong)
	}
	return nil
}

// ComposeFields composes a (string) field list based on a sequence of strings and or
// field lists.
func ComposeFields(fields ...interface{}) Fields {
	var result Fields
	for _, f := range fields {
		switch v := f.(type) {
		case FieldProvider:
			result = append(result, v.GetFields()...)
		case string:
			result = append(result, v)
		case Fields:
			result = append(result, v...)
		case []string:
			result = append(result, v...)
		case []interface{}:
			result = append(result, ComposeFields(v...)...)
		}
	}
	return result
}

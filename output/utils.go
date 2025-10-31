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

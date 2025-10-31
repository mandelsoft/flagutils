package files

import (
	"os"
	"strings"
)

func NameComparator(a, b string) int {
	av := strings.Split(a, string(os.PathSeparator))
	bv := strings.Split(b, string(os.PathSeparator))

	for i, e := range av {
		if i >= len(bv) {
			return 1
		}
		c := strings.Compare(e, bv[i])
		if c != 0 {
			return c
		}
	}
	return len(av) - len(bv)
}

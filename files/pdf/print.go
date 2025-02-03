package pdf

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func PrintMetadata(fields map[string]string) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 2, '\t', 0)

	for field, value := range fields {
		fmt.Fprintf(w, "%s\t%s\n", field, value)
	}

	w.Flush()
}

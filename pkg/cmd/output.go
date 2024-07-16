package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"sigs.k8s.io/yaml"
)

// NewPrinter returns a function that prints in the chose format to writer.
func NewPrinter(format string, w io.Writer) func(any) {
	switch format {

	case "json":
		return func(v any) {
			if b, err := json.Marshal(v); err != nil {
				fmt.Fprintf(w, "%v\n", err)
			} else {
				fmt.Fprintf(w, "%v\n", string(b))
			}
		}

	case "json-pretty":
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return func(v any) { _ = encoder.Encode(v) }

	case "yaml":
		return func(v any) { y, _ := yaml.Marshal(v); fmt.Fprintf(w, "%s", y) }

	default:
		return func(v any) { fmt.Fprintf(w, "%v", v) }
	}
}

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"sigs.k8s.io/yaml"
)

// NewPrinter returns a function that prints in the chose format to writer.
func NewPrinter(format string, w io.Writer) func(any) {
	switch format {

	case "json":
		encoder := json.NewEncoder(w)
		return func(v any) { _ = encoder.Encode(v) }

	case "json-pretty":
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return func(v any) { _ = encoder.Encode(v) }

	case "ndjson":
		encoder := json.NewEncoder(w)
		return func(v any) {
			r := reflect.ValueOf(v)
			switch r.Kind() {
			case reflect.Array, reflect.Slice:
				for i := 0; i < r.Len(); i++ {
					_ = encoder.Encode(r.Index(i).Interface())
				}
			default:
				_ = encoder.Encode(v)
			}
		}

	case "yaml":
		return func(v any) { y, _ := yaml.Marshal(v); fmt.Fprintf(w, "%s", y) }

	default:
		return func(v any) { fmt.Fprintf(w, "%v", v) }
	}
}

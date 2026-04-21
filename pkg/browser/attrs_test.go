package browser

import (
	"testing"
)

func TestAttrs_Attributes(t *testing.T) {
	a := Attrs{"shape": "box", "color": "red"}
	attrs := a.Attributes()
	if len(attrs) != 2 {
		t.Fatalf("got %d attributes, want 2", len(attrs))
	}
	got := map[string]string{}
	for _, attr := range attrs {
		got[attr.Key] = attr.Value
	}
	if got["shape"] != "box" {
		t.Errorf("shape = %q, want %q", got["shape"], "box")
	}
	if got["color"] != "red" {
		t.Errorf("color = %q, want %q", got["color"], "red")
	}
}

func TestAttrs_Attributes_empty(t *testing.T) {
	a := Attrs{}
	attrs := a.Attributes()
	if len(attrs) != 0 {
		t.Errorf("got %d attributes, want 0", len(attrs))
	}
}

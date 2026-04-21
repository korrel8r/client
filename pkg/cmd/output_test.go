package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewPrinter_yaml(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter("yaml", &buf)
	p(map[string]string{"key": "value"})
	got := buf.String()
	if !strings.Contains(got, "key: value") {
		t.Errorf("yaml output = %q, want to contain 'key: value'", got)
	}
}

func TestNewPrinter_json(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter("json", &buf)
	p(map[string]string{"key": "value"})
	var m map[string]string
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if m["key"] != "value" {
		t.Errorf("got %v, want key=value", m)
	}
}

func TestNewPrinter_jsonPretty(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter("json-pretty", &buf)
	p(map[string]string{"key": "value"})
	got := buf.String()
	if !strings.Contains(got, "  ") {
		t.Errorf("json-pretty output not indented: %q", got)
	}
	var m map[string]string
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if m["key"] != "value" {
		t.Errorf("got %v, want key=value", m)
	}
}

func TestNewPrinter_ndjson_slice(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter("ndjson", &buf)
	p([]string{"a", "b", "c"})
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3: %q", len(lines), buf.String())
	}
	for i, want := range []string{`"a"`, `"b"`, `"c"`} {
		got := strings.TrimSpace(lines[i])
		if got != want {
			t.Errorf("line %d = %q, want %q", i, got, want)
		}
	}
}

func TestNewPrinter_ndjson_nonSlice(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter("ndjson", &buf)
	p("scalar")
	var got string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got != "scalar" {
		t.Errorf("got %q, want %q", got, "scalar")
	}
}

func TestNewPrinter_default(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter("unknown-format", &buf)
	p("hello")
	if buf.String() != "hello" {
		t.Errorf("got %q, want %q", buf.String(), "hello")
	}
}

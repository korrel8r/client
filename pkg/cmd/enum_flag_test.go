package cmd

import (
	"strings"
	"testing"
)

func TestEnumFlag_defaults(t *testing.T) {
	f := EnumFlag("a", "b", "c")
	if f.String() != "a" {
		t.Errorf("default = %q, want %q", f.String(), "a")
	}
}

func TestEnumFlag_Set_valid(t *testing.T) {
	f := EnumFlag("a", "b", "c")
	if err := f.Set("b"); err != nil {
		t.Fatalf("Set(b) error: %v", err)
	}
	if f.String() != "b" {
		t.Errorf("after Set(b) = %q, want %q", f.String(), "b")
	}
}

func TestEnumFlag_Set_invalid(t *testing.T) {
	f := EnumFlag("a", "b", "c")
	err := f.Set("invalid")
	if err == nil {
		t.Fatal("Set(invalid) should return error")
	}
	if !strings.Contains(err.Error(), "invalid") {
		t.Errorf("error = %q, should contain 'invalid'", err)
	}
}

func TestEnumFlag_Type(t *testing.T) {
	f := EnumFlag("x", "y")
	got := f.Type()
	if !strings.Contains(got, "x") || !strings.Contains(got, "y") {
		t.Errorf("Type() = %q, should list allowed values", got)
	}
}

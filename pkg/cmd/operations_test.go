package cmd

import (
	"testing"
	"time"
)

func TestBoolPtr_true(t *testing.T) {
	p := boolPtr(true)
	if p == nil || !*p {
		t.Errorf("boolPtr(true) = %v, want *true", p)
	}
}

func TestBoolPtr_false(t *testing.T) {
	p := boolPtr(false)
	if p != nil {
		t.Errorf("boolPtr(false) = %v, want nil", p)
	}
}

func TestGraphOptions_none(t *testing.T) {
	rules, results, errors = false, false, false
	got := graphOptions()
	if got != nil {
		t.Errorf("graphOptions() = %v, want nil", got)
	}
}

func TestGraphOptions_rules(t *testing.T) {
	rules, results, errors = true, false, false
	defer func() { rules, results, errors = false, false, false }()
	got := graphOptions()
	if got == nil {
		t.Fatal("graphOptions() = nil, want non-nil")
	}
	if got.Rules == nil || !*got.Rules {
		t.Error("Rules should be true")
	}
	if got.Results != nil {
		t.Error("Results should be nil")
	}
	if got.Errors != nil {
		t.Error("Errors should be nil")
	}
}

func TestGraphOptions_all(t *testing.T) {
	rules, results, errors = true, true, true
	defer func() { rules, results, errors = false, false, false }()
	got := graphOptions()
	if got == nil {
		t.Fatal("graphOptions() = nil, want non-nil")
	}
	if got.Rules == nil || !*got.Rules {
		t.Error("Rules should be true")
	}
	if got.Results == nil || !*got.Results {
		t.Error("Results should be true")
	}
	if got.Errors == nil || !*got.Errors {
		t.Error("Errors should be true")
	}
}

func TestStart(t *testing.T) {
	class = "mock:mock"
	queries = []string{"mock:mock:{\"a\":\"b\"}"}
	objects = []string{`{"x":1}`}
	defer func() { class = ""; queries = nil; objects = nil }()

	s := start()
	if s.Class != "mock:mock" {
		t.Errorf("Class = %q, want %q", s.Class, "mock:mock")
	}
	if len(s.Queries) != 1 || s.Queries[0] != queries[0] {
		t.Errorf("Queries = %v, want %v", s.Queries, queries)
	}
	if len(s.Objects) != 1 {
		t.Fatalf("Objects len = %d, want 1", len(s.Objects))
	}
}

func TestConstraint_empty(t *testing.T) {
	limit = 0
	since, until = 0, 0
	c := constraint()
	if c.Limit != nil {
		t.Errorf("Limit = %v, want nil", c.Limit)
	}
}

func TestConstraint_withLimit(t *testing.T) {
	limit = 10
	defer func() { limit = 0 }()
	c := constraint()
	if c.Limit == nil || *c.Limit != 10 {
		t.Errorf("Limit = %v, want 10", c.Limit)
	}
}

func TestConstraint_withSinceUntil(t *testing.T) {
	since = 5 * time.Minute
	until = 1 * time.Minute
	defer func() { since, until = 0, 0 }()

	before := time.Now()
	c := constraint()
	after := time.Now()

	if c.Start == nil {
		t.Fatal("Start should not be nil")
	}
	if c.Start.Before(before.Add(-since-time.Second)) || c.Start.After(after.Add(-since+time.Second)) {
		t.Errorf("Start = %v, expected around %v ago", c.Start, since)
	}
	if c.End == nil {
		t.Fatal("End should not be nil")
	}
	if c.End.Before(before.Add(-until-time.Second)) || c.End.After(after.Add(-until+time.Second)) {
		t.Errorf("End = %v, expected around %v ago", c.End, until)
	}
}

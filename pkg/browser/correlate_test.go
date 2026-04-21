package browser

import (
	"net/url"
	"testing"

	"github.com/korrel8r/client/pkg/api"
	"k8s.io/utils/ptr"
)

func TestApiError_success(t *testing.T) {
	for _, code := range []int{200, 201, 204, 299} {
		if err := apiError(code, nil); err != nil {
			t.Errorf("apiError(%d) = %v, want nil", code, err)
		}
	}
}

func TestApiError_withMessage(t *testing.T) {
	body := []byte(`{"error": "not found"}`)
	err := apiError(404, body)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "not found" {
		t.Errorf("error = %q, want %q", err.Error(), "not found")
	}
}

func TestApiError_withoutJSON(t *testing.T) {
	err := apiError(500, []byte("Internal Server Error"))
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "HTTP 500 error" {
		t.Errorf("error = %q, want %q", err.Error(), "HTTP 500 error")
	}
}

func TestApiError_emptyErrorField(t *testing.T) {
	body := []byte(`{"error": ""}`)
	err := apiError(400, body)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "HTTP 400 error" {
		t.Errorf("error = %q, want %q", err.Error(), "HTTP 400 error")
	}
}

func TestCorrelate_addErr_nil(t *testing.T) {
	c := &correlate{}
	if c.addErr(nil) {
		t.Error("addErr(nil) should return false")
	}
	if c.Err != nil {
		t.Errorf("Err = %v, want nil", c.Err)
	}
}

func TestCorrelate_addErr_withError(t *testing.T) {
	c := &correlate{}
	if !c.addErr(errForTest("test error")) {
		t.Error("addErr should return true")
	}
	if c.Err == nil || c.Err.Error() != "test error" {
		t.Errorf("Err = %v, want 'test error'", c.Err)
	}
}

func TestCorrelate_addErr_withPrefix(t *testing.T) {
	c := &correlate{}
	c.addErr(errForTest("detail"), "prefix")
	if c.Err == nil {
		t.Fatal("Err should not be nil")
	}
	got := c.Err.Error()
	if got != "prefix: detail" {
		t.Errorf("Err = %q, want %q", got, "prefix: detail")
	}
}

func TestCorrelate_addErr_accumulates(t *testing.T) {
	c := &correlate{}
	c.addErr(errForTest("first"))
	c.addErr(errForTest("second"))
	if c.Err == nil {
		t.Fatal("Err should not be nil")
	}
	got := c.Err.Error()
	if got != "first\nsecond" {
		t.Errorf("Err = %q, want %q", got, "first\nsecond")
	}
}

func TestCorrelate_NewStartURL(t *testing.T) {
	base, _ := url.Parse("http://localhost/correlate?start=old&goal=3")
	c := &correlate{URL: base}
	got := c.NewStartURL("new-query")
	if got.Query().Get("start") != "new-query" {
		t.Errorf("start = %q, want %q", got.Query().Get("start"), "new-query")
	}
	if got.Query().Get("goal") != "3" {
		t.Errorf("goal = %q, want %q", got.Query().Get("goal"), "3")
	}
}

func TestCorrelate_reset(t *testing.T) {
	b := &Browser{}
	c := &correlate{Browser: b, Err: errForTest("old error")}
	u, _ := url.Parse("http://localhost/correlate?start=q1&goal=log:app")
	c.reset(u)
	if c.Start != "q1" {
		t.Errorf("Start = %q, want %q", c.Start, "q1")
	}
	if c.Goal != "log:app" {
		t.Errorf("Goal = %q, want %q", c.Goal, "log:app")
	}
	if c.Err != nil {
		t.Errorf("Err = %v, want nil", c.Err)
	}
	if c.Browser != b {
		t.Error("Browser should be preserved")
	}
}

func TestCorrelate_reset_defaultGoal(t *testing.T) {
	c := &correlate{Browser: &Browser{}}
	u, _ := url.Parse("http://localhost/correlate?start=q1")
	c.reset(u)
	if c.Goal != "3" {
		t.Errorf("Goal = %q, want default %q", c.Goal, "3")
	}
}

func TestNodeToolTip(t *testing.T) {
	model := &api.Graph{
		Nodes: []api.Node{
			{Class: "k8s:Pod", Queries: []api.QueryCount{
				{Query: "k8s:Pod:{}", Count: ptr.To(5)},
			}},
			{Class: "log:app"},
		},
		Edges: []api.Edge{
			{
				Start: "k8s:Pod",
				Goal:  "log:app",
				Rules: []api.Rule{
					{Name: "rule1", Queries: []api.QueryCount{
						{Query: "k8s:Pod:{}", Count: ptr.To(5)},
					}},
				},
			},
		},
	}
	g := NewGraph(model)
	n := g.NodeFor("k8s:Pod")
	tip := nodeToolTip(g, n)
	if tip == "" {
		t.Error("tooltip should not be empty")
	}
	if len(tip) == 0 {
		t.Error("expected non-empty tooltip")
	}
}

func TestNodeToolTip_noResults(t *testing.T) {
	model := &api.Graph{
		Nodes: []api.Node{
			{Class: "k8s:Pod", Queries: []api.QueryCount{
				{Query: "k8s:Pod:{}", Count: ptr.To(0)},
			}},
		},
	}
	g := NewGraph(model)
	n := g.NodeFor("k8s:Pod")
	tip := nodeToolTip(g, n)
	if tip != "" {
		t.Errorf("tooltip = %q, want empty for zero-count queries", tip)
	}
}

type testErr string

func errForTest(msg string) error { return testErr(msg) }
func (e testErr) Error() string   { return string(e) }

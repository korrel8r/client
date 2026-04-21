package browser

import (
	"testing"

	"github.com/korrel8r/client/pkg/api"
	"k8s.io/utils/ptr"
)

func TestNewGraph_nil(t *testing.T) {
	g := NewGraph(nil)
	if g == nil {
		t.Fatal("NewGraph(nil) returned nil")
	}
	if g.Nodes().Len() != 0 {
		t.Errorf("nodes = %d, want 0", g.Nodes().Len())
	}
	if g.DOTID() != "korrel8r" {
		t.Errorf("DOTID = %q, want %q", g.DOTID(), "korrel8r")
	}
}

func TestNewGraph_withNodes(t *testing.T) {
	model := &api.Graph{
		Nodes: []api.Node{
			{Class: "k8s:Pod", Count: ptr.To(3), Queries: []api.QueryCount{{Query: "k8s:Pod:{}"}}},
			{Class: "log:application", Count: ptr.To(5)},
		},
	}
	g := NewGraph(model)
	if g.Nodes().Len() != 2 {
		t.Fatalf("nodes = %d, want 2", g.Nodes().Len())
	}
	n := g.NodeFor("k8s:Pod")
	if n == nil {
		t.Fatal("NodeFor(k8s:Pod) = nil")
	}
	if n.Model.Class != "k8s:Pod" {
		t.Errorf("class = %q, want %q", n.Model.Class, "k8s:Pod")
	}
}

func TestNewGraph_withEdges(t *testing.T) {
	model := &api.Graph{
		Nodes: []api.Node{
			{Class: "k8s:Pod"},
			{Class: "log:application"},
		},
		Edges: []api.Edge{
			{Start: "k8s:Pod", Goal: "log:application"},
		},
	}
	g := NewGraph(model)
	if g.Edges().Len() != 1 {
		t.Fatalf("edges = %d, want 1", g.Edges().Len())
	}

	from := g.NodeFor("k8s:Pod")
	to := g.NodeFor("log:application")
	if !g.HasEdgeFromTo(from.ID(), to.ID()) {
		t.Error("expected edge from k8s:Pod to log:application")
	}
}

func TestNewGraph_NodeFor_missing(t *testing.T) {
	g := NewGraph(&api.Graph{
		Nodes: []api.Node{{Class: "k8s:Pod"}},
	})
	if n := g.NodeFor("nonexistent"); n != nil {
		t.Errorf("NodeFor(nonexistent) = %v, want nil", n)
	}
}

func TestNewGraph_DOTAttributers(t *testing.T) {
	g := NewGraph(nil)
	ga, na, ea := g.DOTAttributers()
	if ga == nil || na == nil || ea == nil {
		t.Fatal("DOTAttributers returned nil")
	}
	gAttrs := ga.Attributes()
	found := false
	for _, a := range gAttrs {
		if a.Key == "layout" && a.Value == "dot" {
			found = true
		}
	}
	if !found {
		t.Error("expected graph attribute layout=dot")
	}
}

func TestEdge_FromTo(t *testing.T) {
	model := &api.Graph{
		Nodes: []api.Node{
			{Class: "a"},
			{Class: "b"},
		},
		Edges: []api.Edge{
			{Start: "a", Goal: "b"},
		},
	}
	g := NewGraph(model)
	edges := g.Edges()
	if !edges.Next() {
		t.Fatal("no edges")
	}
	e := edges.Edge().(*Edge)
	if e.From().(*Node).Model.Class != "a" {
		t.Errorf("From = %q, want %q", e.From().(*Node).Model.Class, "a")
	}
	if e.To().(*Node).Model.Class != "b" {
		t.Errorf("To = %q, want %q", e.To().(*Node).Model.Class, "b")
	}
}

func TestNode_Attributes(t *testing.T) {
	n := &Node{Model: &api.Node{Class: "test"}, Attrs: Attrs{"shape": "box"}}
	attrs := n.Attributes()
	if len(attrs) != 1 || attrs[0].Key != "shape" {
		t.Errorf("Attributes = %v, want [{shape box}]", attrs)
	}
}

func TestEdge_Attributes(t *testing.T) {
	e := &Edge{
		Edge:  &api.Edge{Start: "a", Goal: "b"},
		Attrs: Attrs{"color": "red"},
		from:  &Node{Model: &api.Node{Class: "a"}},
		to:    &Node{Model: &api.Node{Class: "b"}},
	}
	attrs := e.Attributes()
	if len(attrs) != 1 || attrs[0].Key != "color" {
		t.Errorf("Attributes = %v, want [{color red}]", attrs)
	}
}

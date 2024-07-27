// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package browser

import (
	"unsafe"

	"github.com/korrel8r/client/pkg/swagger/models"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/simple"
)

var (
	_ graph.Node = &Node{}
	_ graph.Edge = &Edge{}
)

type Graph struct {
	Model                            *models.Graph
	GraphAttrs, NodeAttrs, EdgeAttrs Attrs

	nodes map[string]*Node

	*simple.DirectedGraph
}

func NewGraph(mg *models.Graph) *Graph {
	g := &Graph{
		DirectedGraph: simple.NewDirectedGraph(),
		Model:         mg,
		GraphAttrs: Attrs{
			"fontname":        "Helvetica",
			"fontsize":        "12",
			"splines":         "true",
			"overlap":         "prism",
			"overlap_scaling": "-2",
			"layout":          "dot",
		},
		NodeAttrs: Attrs{
			"fontname": "Helveticax",
			"fontsize": "12",
		},
		EdgeAttrs: Attrs{
			"fontname": "Helvetica",
			"fontsize": "12",
		},

		nodes: map[string]*Node{},
	}
	if mg == nil {
		return g
	}
	for _, n := range mg.Nodes {
		nn := &Node{Model: n, Attrs: Attrs{}}
		g.nodes[n.Class] = nn
		g.AddNode(nn)
	}
	for _, e := range mg.Edges {
		g.DirectedGraph.SetEdge(&Edge{
			Edge:  e,
			Attrs: Attrs{},
			from:  g.NodeFor(e.Start),
			to:    g.NodeFor(e.Goal)})
	}
	return g
}

func (g *Graph) DOTID() string { return "korrel8r" }
func (g *Graph) DOTAttributers() (graph, node, edge encoding.Attributer) {
	return g.GraphAttrs, g.NodeAttrs, g.EdgeAttrs
}

type Node struct {
	Model *models.Node
	Attrs
}

func (g *Graph) NodeFor(class string) *Node { return g.nodes[class] }
func (n *Node) ID() int64                   { return id(n) }

type Edge struct {
	*models.Edge
	Attrs
	from, to *Node
}

func (e *Edge) ID() int64                { return id(e) }
func (e *Edge) From() graph.Node         { return e.from }
func (e *Edge) To() graph.Node           { return e.to }
func (e *Edge) ReversedEdge() graph.Edge { panic("not implemented") }

func id[T any](ptr *T) int64 { return (int64)((uintptr)(unsafe.Pointer(ptr))) }

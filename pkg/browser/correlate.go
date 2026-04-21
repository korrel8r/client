// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package browser

import (
	_ "embed"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/korrel8r/client/pkg/api"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"k8s.io/utils/ptr"
)

// correlate web page handler.
type correlate struct {
	URL *url.URL

	// URL Query parameter fields

	Start string // Start query
	Goal  string // Goal class or neighbourhood depth

	// Computed fields used by page template.

	Depth                           int
	Graph                           *Graph
	Diagram, DiagramTxt, DiagramImg string
	ConsoleURL                      *url.URL
	UpdateTime                      time.Duration

	// Other context

	Err     error // Accumulated errors from template.
	Browser *Browser
}

// reset the fields to contain only URL query parameters
func (c *correlate) reset(url *url.URL) {
	params := url.Query()
	app := c.Browser // Save
	*c = correlate{  // Overwrite
		URL:     url,
		Start:   params.Get("start"),
		Goal:    params.Get("goal"),
		Browser: app,
		Graph:   NewGraph(nil),
	}
	if c.Goal == "" {
		c.Goal = "3" // Default to neighbourhood of depth 3
	}
}

func (c *correlate) HTML(gc *gin.Context) {
	c.update(gc.Request)
	if c.Err != nil {
		c.Graph = NewGraph(nil)
	}
	gc.HTML(http.StatusOK, "correlate.html.tmpl", c)
}

func (c *correlate) NewStartURL(query string) *url.URL {
	values := c.URL.Query()
	values.Set("start", query) // Replace start query
	u := url.URL(*c.URL)       // Copy
	u.RawQuery = values.Encode()
	return &u
}

// addErr adds an error to be displayed on the page.
func (c *correlate) addErr(err error, msg ...any) bool {
	if err == nil {
		return false
	}
	switch len(msg) {
	case 0: // Use err unmodified
	case 1: // Use bare msg string as prefix
		err = fmt.Errorf("%v: %w", msg[0], err)
	default: // Treat msg as printf format
		err = fmt.Errorf(msg[0].(string), msg[1:])
	}
	c.Err = errors.Join(c.Err, err)
	return true
}

// apiError extracts an error from a non-2xx response, or returns nil for success.
func apiError(statusCode int, body []byte) error {
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}
	var apiErr api.Error
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error != "" {
		return errors.New(apiErr.Error)
	}
	return fmt.Errorf("HTTP %d error", statusCode)
}

func (c *correlate) update(req *http.Request) {
	c.reset(req.URL)
	start := api.Start{Queries: []string{c.Start}}
	if c.Goal == "" {
		c.addErr(errors.New("search requires a goal class or neighbourhood depth"))
		return
	}
	ctx := context.Background()
	var err error
	c.Depth, err = strconv.Atoi(c.Goal)
	if err == nil {
		params := &api.GraphNeighborsParams{
			Options: &api.GraphOptions{Rules: ptr.To(true)},
		}
		body := api.Neighbors{
			Start: start,
			Depth: c.Depth,
		}
		resp, respErr := c.Browser.client.GraphNeighborsWithResponse(ctx, params, body)
		if respErr != nil {
			c.addErr(respErr)
		} else if apiErr := apiError(resp.StatusCode(), resp.Body); apiErr != nil {
			c.addErr(apiErr)
		} else {
			c.Graph = NewGraph(resp.JSON200)
		}
	} else {
		params := &api.GraphGoalsParams{
			Options: &api.GraphOptions{Rules: ptr.To(true)},
		}
		body := api.Goals{
			Start: start,
			Goals: []string{c.Goal},
		}
		resp, respErr := c.Browser.client.GraphGoalsWithResponse(ctx, params, body)
		if respErr != nil {
			c.addErr(respErr)
		} else if apiErr := apiError(resp.StatusCode(), resp.Body); apiErr != nil {
			c.addErr(apiErr)
		} else {
			c.Graph = NewGraph(resp.JSON200)
		}
	}
	c.updateDiagram()
}

var domainAttrs = map[string]Attrs{
	"k8s":     {"shape": "septagon", "fillcolor": "#326CE5", "fontcolor": "white"},
	"log":     {"shape": "note", "fillcolor": "yellow"},
	"alert":   {"shape": "triangle", "fillcolor": "yellow"},
	"metric":  {"shape": "egg", "fillcolor": "wheat"},
	"netflow": {"shape": "component", "fillcolor": "skyblue"},
	"trace":   {"shape": "folder", "fillcolor": "aquamarine"},
}

func nodeToolTip(g *Graph, n *Node) string {
	// Collect rules that contributed to each query on node.
	rules := map[string][]string{}
	edges := g.Edges()
	for edges.Next() {
		e := edges.Edge().(*Edge)
		for _, r := range e.Rules {
			for _, q := range r.Queries {
				if q.Count != nil && *q.Count > 0 {
					rules[q.Query] = append(rules[q.Query], r.Name)
				}
			}
		}
	}
	// Build tool tip text
	w := &strings.Builder{}
	for _, q := range n.Model.Queries {
		if q.Count != nil && *q.Count > 0 {
			fmt.Fprintf(w, "%v %v %v\n", *q.Count, rules[q.Query], q.Query)
		}
	}
	return w.String()
}

// updateDiagram generates an SVG diagram via graphviz.
func (c *correlate) updateDiagram() {
	nodes := c.Graph.Nodes()
	for nodes.Next() {
		n := nodes.Node().(*Node)
		var count int
		if n.Model.Count != nil {
			count = *n.Model.Count
		}
		a := n.Attrs
		a["style"] = "filled"
		a["label"] = fmt.Sprintf("%v\n%v", n.Model.Class, count)
		a["tooltip"] = nodeToolTip(c.Graph, n)
		maps.Copy(a, domainAttrs[strings.SplitN(n.Model.Class, ":", 2)[0]])
	}
	// Write the graph files
	baseName := filepath.Join(c.Browser.dir, "files", "korrel8r")
	if gv, err := dot.Marshal(c.Graph, "", "", "  "); !c.addErr(err) {
		gvFile := baseName + ".txt"
		if !c.addErr(os.WriteFile(gvFile, gv, 0664)) {
			// Render and write the graph image
			svgFile := baseName + ".svg"
			if !c.addErr(runDot("dot", "-v", "-Tsvg", "-o", svgFile, gvFile)) {
				c.Diagram, _ = filepath.Rel(c.Browser.dir, svgFile)
				c.DiagramTxt, _ = filepath.Rel(c.Browser.dir, gvFile)
				pngFile := baseName + ".png"
				if !c.addErr(runDot("dot", "-v", "-Tpng", "-o", pngFile, gvFile)) {
					c.DiagramImg, _ = filepath.Rel(c.Browser.dir, pngFile)
				}
			}
		}
	}
}

func runDot(cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v %w: %v", cmdName, err, string(out))
	}
	return err
}

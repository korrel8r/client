// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package cmd_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/korrel8r/client/pkg/api"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

var allDomains = []string{"alert", "incident", "k8s", "log", "metric", "netflow", "mock", "trace"}

func Test_version(t *testing.T) {
	out, err := korrel8rcli(t, "version")
	require.NoError(t, err)
	require.NotEmpty(t, strings.TrimSpace(out))
}

func Test_domains(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "domains", "-u", u.String())
	require.NoError(t, err)

	var domains []api.Domain
	require.NoError(t, yaml.Unmarshal([]byte(out), &domains), out, string(out))
	var names []string
	for _, d := range domains {
		names = append(names, d.Name)
	}
	require.ElementsMatch(t, allDomains, names)
}

func Test_domains_json(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "domains", "-u", u.String(), "-o", "json")
	require.NoError(t, err)
	var domains []api.Domain
	require.NoError(t, json.Unmarshal([]byte(out), &domains))
	require.Len(t, domains, len(allDomains))
}

func Test_domains_jsonPretty(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "domains", "-u", u.String(), "-o", "json-pretty")
	require.NoError(t, err)
	require.Contains(t, out, "\n  ")
	var domains []api.Domain
	require.NoError(t, json.Unmarshal([]byte(out), &domains))
	require.Len(t, domains, len(allDomains))
}

func Test_domains_ndjson(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "domains", "-u", u.String(), "-o", "ndjson")
	require.NoError(t, err)
	// ndjson with a pointer-to-slice prints the whole array as one JSON line
	var domains []api.Domain
	require.NoError(t, json.Unmarshal([]byte(out), &domains))
	require.Len(t, domains, len(allDomains))
}

func Test_classes(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "classes", "-u", u.String(), "alert")
	require.NoError(t, err)
	var classes []string
	require.NoError(t, yaml.Unmarshal([]byte(out), &classes))
	require.Contains(t, classes, "alert")
}

func Test_classes_k8s(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "classes", "-u", u.String(), "k8s")
	require.NoError(t, err)
	var classes []string
	require.NoError(t, yaml.Unmarshal([]byte(out), &classes))
	require.Contains(t, classes, "Event.v1")
}

func Test_classes_invalidDomain(t *testing.T) {
	u := korrel8rServer(t)
	_, err := korrel8rcli(t, "classes", "-u", u.String(), "nosuchdomain")
	require.Error(t, err)
	require.Contains(t, err.Error(), "nosuchdomain")
}

func Test_classes_missingArg(t *testing.T) {
	u := korrel8rServer(t)
	_, err := korrel8rcli(t, "classes", "-u", u.String())
	require.Error(t, err)
}

func Test_objects_invalidQuery(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "objects", "-u", u.String(), "this-is-not-a-query")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid query: this-is-not-a-query")
	require.Equal(t, "", out)
}

func Test_objects_noStore(t *testing.T) {
	u := korrel8rServer(t)
	_, err := korrel8rcli(t, "objects", "-u", u.String(), "alert:alert:{}")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no stores found")
}

func Test_neighbors(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "neighbors", "-u", u.String(), "-q", "log:application:{}", "-d", "1")
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func Test_neighbors_depth2(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "neighbors", "-u", u.String(), "-q", "log:application:{}", "-d", "2")
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func Test_neighbors_withRules(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "neighbors", "-u", u.String(), "-q", "log:application:{}", "-d", "1", "--rules", "-o", "json")
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func Test_neighbors_withAllFlags(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "neighbors", "-u", u.String(),
		"-q", "log:application:{}",
		"-d", "1",
		"--rules", "--results", "--errors",
		"-o", "json")
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func Test_neighbors_invalidClass(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "neighbors", "-u", u.String(), "--class", "invalid:class:name", "-d", "1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "POST /graphs/neighbors")
	require.Equal(t, "", out)
}

func Test_goals(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "goals", "-u", u.String(), "-q", "log:application:{}", "k8s:Event.v1")
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func Test_goals_withRules(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "goals", "-u", u.String(),
		"-q", "log:application:{}",
		"--rules",
		"-o", "json",
		"k8s:Event.v1")
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func Test_goals_multipleGoals(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "goals", "-u", u.String(),
		"-q", "log:application:{}",
		"k8s:Event.v1", "metric:metric")
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func Test_goals_missingArgs(t *testing.T) {
	u := korrel8rServer(t)
	_, err := korrel8rcli(t, "goals", "-u", u.String())
	require.Error(t, err)
}

func Test_listGoals(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "list-goals", "-u", u.String(), "-q", "log:application:{}", "k8s:Event.v1")
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func Test_listGoals_json(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "list-goals", "-u", u.String(), "-q", "log:application:{}", "-o", "json", "k8s:Event.v1")
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func Test_config(t *testing.T) {
	u := korrel8rServer(t)
	_, err := korrel8rcli(t, "config", "-u", u.String())
	require.NoError(t, err)
}

func Test_config_setVerbose(t *testing.T) {
	u := korrel8rServer(t)
	_, err := korrel8rcli(t, "config", "-u", u.String(), "--set-verbose", "1")
	require.NoError(t, err)
}

func Test_setConsole(t *testing.T) {
	u := korrel8rServer(t)
	_, err := korrel8rcli(t, "set-console", "-u", u.String(), `{"view":"k8s:Pod.v1:{}" }`)
	require.NoError(t, err)
}

func Test_setConsole_invalidJSON(t *testing.T) {
	u := korrel8rServer(t)
	_, err := korrel8rcli(t, "set-console", "-u", u.String(), "not-json")
	require.Error(t, err)
}

func Test_setConsole_missingArg(t *testing.T) {
	u := korrel8rServer(t)
	_, err := korrel8rcli(t, "set-console", "-u", u.String())
	require.Error(t, err)
}

func Test_invalidOutputFormat(t *testing.T) {
	u := korrel8rServer(t)
	_, err := korrel8rcli(t, "domains", "-u", u.String(), "-o", "invalid")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid")
}

func Test_error_includes_http_context(t *testing.T) {
	u := korrel8rServer(t)

	// Test with invalid query - should show GET /objects with error message
	out, err := korrel8rcli(t, "objects", "-u", u.String(), "invalid-query")
	require.Error(t, err)
	require.Contains(t, err.Error(), "GET /objects")
	require.Contains(t, err.Error(), "invalid query")
	require.Equal(t, "", out)

	// Test with invalid class in neighbors - should show POST /graphs/neighbors
	out, err = korrel8rcli(t, "neighbors", "-u", u.String(), "--class", "invalid:class:name", "-d", "1")
	require.Error(t, err)
	// The error should contain HTTP method and endpoint context
	require.Contains(t, err.Error(), "POST /graphs/neighbors", "error should include HTTP method and endpoint")
	require.Equal(t, "", out)
}

func Test_neighbors_withConstraints(t *testing.T) {
	u := korrel8rServer(t)
	out, err := korrel8rcli(t, "neighbors", "-u", u.String(),
		"-q", "log:application:{}",
		"-d", "1",
		"--limit", "10",
		"--since", "1h",
		"--until", "5m")
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

var buildOnce sync.Once

// korrel8rcli returns an exec.Cmd to run the executable in the context of a testing.T test.
// Includes support for writing coverage data to
func korrel8rcli(t *testing.T, args ...string) (out string, err error) {
	t.Helper()
	const (
		dir = "../../cmd/korrel8rcli"
		exe = "../../korrel8rcli"
	)
	buildOnce.Do(func() {
		cmd := exec.Command("go", "build", "-cover", "-o", exe, dir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		require.NoError(t, cmd.Run())
	})
	require.NoError(t, os.MkdirAll("_covdata", 0770))
	cmd := exec.Command(exe, args...)
	cmd.Env = []string{"GOCOVERDIR=_covdata"}
	b, err := cmd.Output()
	if exitErr, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("%w: stderr: %v", exitErr, string(exitErr.Stderr))
	}
	return string(b), err
}

// Start a korrel8r server, will shut down at end of test.
func korrel8r(t *testing.T) *url.URL {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	u := &url.URL{Scheme: "http", Host: l.Addr().String()}
	require.NoError(t, l.Close())
	korrel8rCmd, err := exec.Command("go", "tool", "-n", "korrel8r").Output()
	require.NoError(t, err, "korrel8r must be a go tool dependency")
	cmd := exec.Command(strings.TrimSpace(string(korrel8rCmd)), "web", "--http", u.Host, "-c=testdata/korrel8r.yaml")
	cmd.Stderr = &testWriter{Name: "korrel8r", T: t}
	require.NoError(t, cmd.Start())
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	return u
}

type testWriter struct {
	T    *testing.T
	Name string
}

func (w *testWriter) Write(data []byte) (int, error) {
	w.T.Logf("%v:%v", w.Name, string(data))
	return len(data), nil
}

func korrel8rServer(t *testing.T) *url.URL {
	t.Helper()
	u := korrel8r(t) // Start the server
	var err error
	// Wait for server to be listening
	require.Eventually(t, func() bool {
		c, err := net.Dial("tcp", u.Host)
		if err == nil {
			_ = c.Close()
		}
		return err == nil
	}, time.Second, time.Second/10)
	require.NoError(t, err)
	return u
}

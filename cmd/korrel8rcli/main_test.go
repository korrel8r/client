// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package main_test

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/korrel8r/client/pkg/swagger/models"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func Test_domains(t *testing.T) {
	u := setup(t) // Start the server
	out, err := korrel8rcli(t, "domains", "-u", u.String())
	require.NoError(t, err)

	var domains []*models.Domain
	require.NoError(t, yaml.Unmarshal([]byte(out), &domains))
	var names []string
	for _, d := range domains {
		names = append(names, d.Name)
	}
	require.ElementsMatch(t, []string{"k8s", "alert", "log", "metric", "netflow", "mock"}, names)
}

func Test_bad_parameters(t *testing.T) {
	u := setup(t) // Start the server
	out, err := korrel8rcli(t, "objects", "-u", u.String(), "this-is-not-a-query")
	require.EqualError(t, err, "exit status 1: stderr: korrel8rcli: [GET /objects][400] GetObjects default {\"error\":\"invalid query string: this-is-not-a-query\"}\n")
	require.Equal(t, out, "")
}

// korrel8rcli returns an exec.Cmd to run the executable in the context of a testing.T test.
// Includes support for writing coverage data to
func korrel8rcli(t *testing.T, args ...string) (out string, err error) {
	t.Helper()
	require.NoError(t, os.MkdirAll("covdata", 0770))
	cmd := exec.Command("./korrel8rcli", args...)
	cmd.Env = []string{"GOCOVERDIR=covdata"}
	b, err := cmd.Output()
	if exitErr, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("%w: stderr: %v", exitErr, string(exitErr.Stderr))
	}
	return string(b), err
}

// Start a korrel8r server, will shut down at end of test.
func korrel8r(t *testing.T, args ...string) *url.URL {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	u := &url.URL{Scheme: "http", Host: l.Addr().String()}
	require.NoError(t, l.Close())
	korrel8rCmd := os.Getenv("KORREL8R")
	if korrel8rCmd == "" {
		t.Fatal("Set env var KORREL8R to the korrel8r executable.")
	}
	cmd := exec.Command(korrel8rCmd, "web", "--http", u.Host, "-c=testdata/korrel8r.yaml")
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

func setup(t *testing.T) *url.URL {
	t.Helper()
	u := korrel8r(t) // Start the server
	var err error
	// Wait for server to be listening
	require.Eventually(t, func() bool {
		c, err := net.Dial("tcp", u.Host)
		defer func() { _ = c.Close() }()
		return err == nil
	}, time.Second, time.Second/10)
	require.NoError(t, err)
	return u
}

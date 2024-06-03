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
	u := korrel8r(t) // Start the server
	var (
		out string
		err error
	)
	// Wait for server to start
	require.Eventually(t, func() bool {
		out, err = korrel8rcli(t, "domains", "-u", u.String())
		if err != nil {
			t.Log("retry: ", err)
		}
		return err == nil
	}, time.Second, time.Second/10)
	require.NoError(t, err)

	var domains []*models.Domain
	require.NoError(t, yaml.Unmarshal([]byte(out), &domains))
	var names []string
	for _, d := range domains {
		names = append(names, d.Name)
	}
	require.ElementsMatch(t, []string{"k8s", "alert", "log", "metric", "netflow", "mock"}, names)
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
	cmd := exec.Command(os.Getenv("KORREL8R"), "web", "--http", u.Host, "-c=testdata/korrel8r.yaml")
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

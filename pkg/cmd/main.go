// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package cmd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/korrel8r/client/pkg/api"
	"github.com/korrel8r/client/pkg/build"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const DefaultBasePath = "/api/v1alpha1"

func Main() {
	rootCmd.PersistentFlags().VarP(output, "output", "o", "Output format")
	log.SetPrefix(filepath.Base(os.Args[0]) + ": ")
	log.SetFlags(0)
	check(rootCmd.Execute())
}

const (
	envURL         = "KORREL8RCLI_URL"
	envBearerToken = "KORREL8RCLI_BEARER_TOKEN"
)

var (
	rootCmd = &cobra.Command{
		Use:     "korrel8rcli COMMAND",
		Short:   "REST client for a remote korrel8r server.",
		Version: build.Version,
	}

	// Global Flags
	output      = EnumFlag("yaml", "json-pretty", "json", "ndjson")
	korrel8rURL = rootCmd.PersistentFlags().StringP("url", "u", urlDefault(),
		fmt.Sprintf("URL of remote korrel8r, default from env %v", envURL))
	insecure = rootCmd.PersistentFlags().BoolP("insecure", "k", false, "Insecure connection, skip TLS verification.")
	// NOTE don't show the bearer token default for security reasons.
	bearerTokenFlag = rootCmd.PersistentFlags().StringP("bearer-token", "t", "",
		fmt.Sprintf("Authhorization token, default from env %v or kube config.", envBearerToken))
	debug = rootCmd.PersistentFlags().Bool("debug", false, "Enable debug output.")
)

func urlDefault() string {
	if u := os.Getenv(envURL); u != "" {
		return u
	}
	return "http://localhost:8080"
}
func bearerToken() string {
	if *bearerTokenFlag != "" { // Flag first
		return *bearerTokenFlag
	}
	if token := os.Getenv(envBearerToken); token != "" { // Env next
		return token
	}
	if cfg, err := config.GetConfig(); err == nil { // Kube config last
		if cfg.BearerTokenFile != "" { // Try the file first
			if b, err := os.ReadFile(cfg.BearerTokenFile); err == nil {
				return string(b)
			}
		}
		return cfg.BearerToken
	}
	return ""
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version.",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println(rootCmd.Version) },
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func newClient() *api.ClientWithResponses {
	if *korrel8rURL == "" {
		check(fmt.Errorf("no URL: set --url flag or %v environment variable", envURL))
	}
	u, err := url.Parse(*korrel8rURL)
	check(err)
	if u.Path == "" || u.Path == "/" {
		u.Path = DefaultBasePath
	}

	var opts []api.ClientOption

	// Configure HTTP client with optional TLS
	if *insecure {
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // user-requested insecure mode
			},
		}
		opts = append(opts, api.WithHTTPClient(httpClient))
	}

	// Add bearer token authentication
	if token := bearerToken(); token != "" {
		opts = append(opts, api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+token)
			return nil
		}))
	}

	// Add debug logging
	if *debug {
		opts = append(opts, api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			log.Printf("DEBUG: %s %s", req.Method, req.URL)
			return nil
		}))
	}

	client, err := api.NewClientWithResponses(u.String(), opts...)
	check(err)
	return client
}

func check(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// checkResponse checks an HTTP response for errors and exits with a formatted message if non-2xx.
// Returns the response body for successful responses.
func checkResponse(statusCode int, body []byte, method, path string) {
	if statusCode >= 200 && statusCode < 300 {
		return
	}
	var apiErr api.Error
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error != "" {
		fmt.Fprintf(os.Stderr, "%s %s: %s\n", method, path, apiErr.Error)
	} else {
		fmt.Fprintf(os.Stderr, "%s %s: HTTP %d error\n", method, path, statusCode)
	}
	os.Exit(1)
}

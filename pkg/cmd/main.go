// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/korrel8r/client/pkg/build"
	"github.com/korrel8r/client/pkg/swagger/client"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

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
		fmt.Sprintf("URL of remote korrel8r. Default from env %v", envURL))
	insecure = rootCmd.PersistentFlags().BoolP("insecure", "k", false, "Insecure connection, skip TLS verification.")
	// NOTE don't show the bearer token default for security reasons.
	bearerTokenFlag = rootCmd.PersistentFlags().StringP("bearer-token", "t", "",
		fmt.Sprintf("Authhorization token. Default from %v or kube config.", envBearerToken))
	debug = rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug output.")
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

func newClient() *client.RESTAPI {
	if *korrel8rURL == "" {
		check(fmt.Errorf("Either command line flag --url or environment variable %v be set.", envURL))
	}
	u, err := url.Parse(*korrel8rURL)
	check(err)
	if u.Path == "" || u.Path == "/" {
		u.Path = client.DefaultBasePath
	}
	var transport *httptransport.Runtime
	if *insecure {
		tlsClient, err := httptransport.TLSClient(httptransport.TLSClientOptions{InsecureSkipVerify: *insecure})
		check(err)
		transport = httptransport.NewWithClient(u.Host, u.Path, []string{u.Scheme}, tlsClient)
	} else {
		transport = httptransport.New(u.Host, u.Path, []string{u.Scheme})
	}
	if token := bearerToken(); token != "" {
		transport.DefaultAuthentication = httptransport.BearerToken(token)
	}
	transport.Debug = *debug
	return client.New(transport, nil)
}

func check(err error) {
	if err == nil {
		return
	}
	if hasPayload, ok := err.(interface{ GetPayload() any }); ok {
		if m, ok := hasPayload.GetPayload().(map[string]any); ok {
			if msg := m["error"]; msg != "" {
				fmt.Fprintln(os.Stderr, msg)
				os.Exit(1)
			}
		}
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

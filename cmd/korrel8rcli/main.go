// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"os"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/korrel8r/client/pkg/build"
	"github.com/korrel8r/client/pkg/swagger/client"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "korrel8rcli COMMAND",
		Short:   "REST client for a remote korrel8r server.",
		Version: build.Version,
	}

	// Global Flags
	output          = EnumFlag("yaml", "json-pretty", "json")
	korrel8rURL     = rootCmd.PersistentFlags().StringP("url", "u", os.Getenv("KORREL8RCLI_URL"), "URL of remote korrel8r service, environment KORREL8RCLI_URL")
	insecure        = rootCmd.PersistentFlags().BoolP("insecure", "k", false, "Insecure connection, skip TLS verification.")
	bearerTokenFlag = rootCmd.PersistentFlags().String("bearer-token", "", "Bearer token for Authorization, environment KORREL8RCLI_BEARER_TOKEN")
)

func bearerToken() string {
	if *bearerTokenFlag != "" {
		return *bearerTokenFlag
	}
	return os.Getenv("KORREL8RCLI_BEARER_TOKEN")
}

func main() {
	rootCmd.PersistentFlags().VarP(output, "output", "o", "Output format")
	log.SetPrefix(filepath.Base(os.Args[0]) + ": ")
	log.SetFlags(0)
	check(rootCmd.Execute())
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
		check(errors.New("Either command line flag --url or environment variable KORREL8R_URL must be set. "))
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
	if bearerToken() != "" {
		transport.DefaultAuthentication = httptransport.BearerToken(bearerToken())
	}
	return client.New(transport, nil)
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

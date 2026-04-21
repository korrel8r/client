// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/korrel8r/client/pkg/api"
	"github.com/spf13/cobra"
	"k8s.io/utils/ptr"
)

// Common flags for neighbours and goals
var (
	class   string
	queries []string
	objects []string
	rules   bool
	results bool
	errors  bool

	limit        int
	since, until time.Duration
)

func startFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&queries, "query", "q", nil, "Query string for start objects, can be multiple.")
	cmd.Flags().StringVarP(&class, "class", "c", "", "Class for serialized start objects")
	cmd.Flags().StringArrayVarP(&objects, "object", "O", nil, "Serialized start object, can be multiple.")
	cmd.Flags().IntVar(&limit, "limit", 0, "Limit total number of results.")
	cmd.Flags().DurationVar(&since, "since", 0, "Only get results since this long ago.")
	cmd.Flags().DurationVar(&until, "until", 0, "Only get results until this long ago.")
}

func graphFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&rules, "rules", false, "Include rule information in returned graph.")
	cmd.Flags().BoolVar(&results, "results", false, "Include full JSON results with each query.")
	cmd.Flags().BoolVar(&errors, "errors", false, "Include non-fatal error messages.")
}

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Get a list of domains and store configuration",
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		resp, err := c.ListDomainsWithResponse(context.Background())
		check(err)
		checkResponse(resp.StatusCode(), resp.Body, "GET", "/domains")
		NewPrinter(output.String(), os.Stdout)(resp.JSON200)
	},
}

func init() {
	rootCmd.AddCommand(domainsCmd)
}

var classesCmd = &cobra.Command{
	Use:   "classes DOMAIN",
	Short: "Get the list of classes for a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		resp, err := c.ListDomainClassesWithResponse(context.Background(), args[0])
		check(err)
		checkResponse(resp.StatusCode(), resp.Body, "GET", "/domain/"+args[0]+"/classes")
		NewPrinter(output.String(), os.Stdout)(resp.JSON200)
	},
}

func init() {
	rootCmd.AddCommand(classesCmd)
}

var objectsCmd = &cobra.Command{
	Use:   "objects QUERY",
	Short: "Return the list of objects for a query.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		resp, err := c.ObjectsWithResponse(context.Background(), &api.ObjectsParams{Query: args[0]})
		check(err)
		checkResponse(resp.StatusCode(), resp.Body, "GET", "/objects")
		NewPrinter(output.String(), os.Stdout)(resp.JSON200)
	},
}

func init() {
	rootCmd.AddCommand(objectsCmd)
}

var neighboursCmd = &cobra.Command{
	Use:   "neighbours",
	Short: "Get graph of nearest neighbours",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		params := &api.GraphNeighborsParams{
			Options: graphOptions(),
		}
		body := api.Neighbors{
			Depth: depth,
			Start: start(),
		}
		resp, err := c.GraphNeighborsWithResponse(context.Background(), params, body)
		check(err)
		checkResponse(resp.StatusCode(), resp.Body, "POST", "/graphs/neighbors")
		NewPrinter(output.String(), os.Stdout)(resp.JSON200)
	},
}

var depth int

func init() {
	rootCmd.AddCommand(neighboursCmd)
	startFlags(neighboursCmd)
	graphFlags(neighboursCmd)
	neighboursCmd.Flags().IntVarP(&depth, "depth", "d", 2, "Depth of neighbourhood search.")
}

var goalsCmd = &cobra.Command{
	Use:   "goals CLASS...",
	Short: "Get graph of goal classes reachable from start",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		params := &api.GraphGoalsParams{
			Options: graphOptions(),
		}
		body := api.Goals{
			Goals: args,
			Start: start(),
		}
		resp, err := c.GraphGoalsWithResponse(context.Background(), params, body)
		check(err)
		checkResponse(resp.StatusCode(), resp.Body, "POST", "/graphs/goals")
		NewPrinter(output.String(), os.Stdout)(resp.JSON200)
	},
}

func init() {
	rootCmd.AddCommand(goalsCmd)
	startFlags(goalsCmd)
	graphFlags(goalsCmd)
}

var listGoalsCmd = &cobra.Command{
	Use:   "list-goals CLASS...",
	Short: "List goal nodes related to a starting point",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		body := api.Goals{
			Goals: args,
			Start: start(),
		}
		resp, err := c.ListGoalsWithResponse(context.Background(), body)
		check(err)
		checkResponse(resp.StatusCode(), resp.Body, "POST", "/lists/goals")
		NewPrinter(output.String(), os.Stdout)(resp.JSON200)
	},
}

func init() {
	rootCmd.AddCommand(listGoalsCmd)
	startFlags(listGoalsCmd)
}

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Change configuration settings on the server",
		Run: func(cmd *cobra.Command, args []string) {
			params := &api.SetConfigParams{}
			if cmd.Flags().Changed("set-verbose") {
				params.Verbose = &configVerbose
			}
			c := newClient()
			resp, err := c.SetConfigWithResponse(context.Background(), params)
			check(err)
			checkResponse(resp.StatusCode(), resp.Body, "PUT", "/config")
		},
	}

	configVerbose int
)

func init() {
	configCmd.Flags().IntVar(&configVerbose, "set-verbose", 0, "Set verbose level for logging")
	rootCmd.AddCommand(configCmd)
}

var setConsoleCmd = &cobra.Command{
	Use:   "set-console JSON",
	Short: "Set console state for an agent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var console api.Console
		check(json.Unmarshal([]byte(args[0]), &console))
		c := newClient()
		resp, err := c.SetConsoleWithResponse(context.Background(), console)
		check(err)
		checkResponse(resp.StatusCode(), resp.Body, "PUT", "/console")
	},
}

func init() {
	rootCmd.AddCommand(setConsoleCmd)
}

var consoleEventsCmd = &cobra.Command{
	Use:   "console-events",
	Short: "Stream console events from an agent",
	Long:  "Subscribe to SSE event stream of console display updates from an agent.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		resp, err := c.ConsoleEvents(context.Background())
		check(err)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			checkResponse(resp.StatusCode, body, "GET", "/console/events")
		}
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if after, ok := strings.CutPrefix(line, "data: "); ok {
				fmt.Println(after)
			}
		}
		check(scanner.Err())
	},
}

func init() {
	rootCmd.AddCommand(consoleEventsCmd)
}

func boolPtr(v bool) *bool {
	if v {
		return ptr.To(true)
	}
	return nil
}

func graphOptions() *api.GraphOptions {
	if !rules && !results && !errors {
		return nil
	}
	return &api.GraphOptions{Rules: boolPtr(rules), Results: boolPtr(results), Errors: boolPtr(errors)}
}

func start() api.Start {
	var objs []api.Object
	for _, o := range objects {
		objs = append(objs, json.RawMessage(o))
	}
	return api.Start{
		Class:      class,
		Constraint: constraint(),
		Objects:    objs,
		Queries:    queries,
	}
}

func constraint() *api.Constraint {
	c := &api.Constraint{}
	if limit > 0 {
		c.Limit = ptr.To(limit)
	}
	now := time.Now()
	if since > 0 {
		c.Start = ptr.To(now.Add(-since))
	}
	if until > 0 {
		c.End = ptr.To(now.Add(-until))
	}
	return c
}
